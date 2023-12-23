package main

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/router"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app/usecase/save_race_time"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app/usecase/test_time_setter"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/config"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/config/flags"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/kafka"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/pgqueue"
)

func main() {
	ctx := context.Background()
	flags.InitServiceFlags()

	//- service connections
	db, err := config.ConnectPostgres(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to postgresql ", err)
	}
	logrus.Info("connected to postgres!")

	kafkaProducerClient, err := config.ConnectProducer(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to kafka ", err)
	}

	// common
	var ()

	// interfases
	var (
		kafkaProducer  = kafka.NewKafkaProducer("test_topic", kafkaProducerClient)
		pgqueueClient  = pgqueue.NewPgQueueWorker(db, 100*time.Millisecond, time.Minute, 50, 5)
		testTimeSetter = test_time_setter.NewTestTimeSetter(kafkaProducer)
	)

	// registration
	pgqueueClient.RegisterTask(
		pgqueue.TimeSetKind,
		testTimeSetter,
	)

	pgqueueClient.RegisterTask(
		pgqueue.RaceTime,
		save_race_time.NewRaceSetter(kafkaProducer),
	)

	// Имплементация API
	err = pgqueueClient.Run()
	if err != nil {
		logrus.Fatal(err)
	}
	defer pgqueueClient.Stop()
	application, err := app.NewImplementation(pgqueueClient)
	if err != nil {
		logrus.Fatalln("cannot configure implementation")
	}

	root := router.NewRouter(middleware.NewErrorHandler(application))

	logrus.Info("app successfully started")
	err = http.ListenAndServe(":"+config.GetValue(config.ListenPort), root)
	if err != nil {
		logrus.Fatalln("unexpected error from app")
	}
}
