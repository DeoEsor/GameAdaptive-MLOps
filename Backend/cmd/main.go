package main

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/middleware"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/api/router"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/app"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/config"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/config/flags"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/pkg/kafka"
)

func main() {
	ctx := context.Background()
	flags.InitServiceFlags()

	//- service connections
	_, err := config.ConnectPostgres(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to postgresql ", err)
	}

	kafkaProducerClient, err := config.ConnectProducer(ctx)
	if err != nil {
		logrus.Fatalln("cannot connect to kafka ", err)
	}

	// common
	kafkaProducer := kafka.NewKafkaProducer("test_topic", kafkaProducerClient)

	// interfaces

	// Имплементация API
	application, err := app.NewImplementation(kafkaProducer)
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
