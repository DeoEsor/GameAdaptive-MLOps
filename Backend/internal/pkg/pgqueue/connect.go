package pgqueue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/gen/ml_ops/public/model"
	"github.com/DeoEsor/GameAdaptive-MLOp/Backend/internal/gen/ml_ops/public/table"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	TaskKind    string
	TaskHandler interface {
		Handle(ctx context.Context, taskName TaskKind, payload []byte) error
	}
)

type PgQueue struct {
	tasks map[TaskKind]TaskHandler

	chanClose            chan struct{}
	taskTickDuration     time.Duration
	taskRunAfterDuration time.Duration
	ticker               time.Ticker
	workLimits           int
	attemptsLeftDefault  int

	db *pgxpool.Pool
}

func NewPgQueueWorker(db *pgxpool.Pool, taskTickDuration time.Duration, taskRunAfterDuration time.Duration, workLimits, attemptsLeft int) *PgQueue {
	result := PgQueue{
		tasks:                make(map[TaskKind]TaskHandler),
		taskTickDuration:     taskTickDuration,
		taskRunAfterDuration: taskRunAfterDuration,
		workLimits:           workLimits,
		attemptsLeftDefault:  attemptsLeft,
		db:                   db,
	}
	return &result
}

func (pgq *PgQueue) Run() error {
	if len(pgq.tasks) == 0 {
		return fmt.Errorf("no tasks to run")
	}
	pgq.ticker = *time.NewTicker(pgq.taskTickDuration)
	logrus.Info("starting pgqueue worker...")
	go func() {
		for {
			select {
			case <-pgq.chanClose:
				logrus.Info("pgqueue worker is stopping...")
				return
			case <-pgq.ticker.C:
				attempted, err := pgq.performTask()
				if err != nil {
					logrus.Error(err)
					time.Sleep(time.Second)
				}
				if !attempted {
					time.Sleep(time.Second)
				}
			}
		}
	}()
	return nil
}

func (pgq *PgQueue) RegisterTask(taskName TaskKind, taskHandler TaskHandler) {
	pgq.tasks[taskName] = taskHandler
}

func (pgq *PgQueue) Stop() {
	pgq.ticker.Stop()
	pgq.chanClose <- struct{}{}
}

func (pgq *PgQueue) performTask() (attempted bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := pgq.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("cannot create tx: %w", err)
	}
	defer tx.Rollback(ctx)

	chooseTasksQuery, args := table.PgqJobs.SELECT(
		table.PgqJobs.ID,
		table.PgqJobs.CreatedAt,
		table.PgqJobs.TaskName,
		table.PgqJobs.Payload,
		table.PgqJobs.RunAfter,
		table.PgqJobs.RetryWaits,
		table.PgqJobs.InRun,
	).WHERE(
		table.PgqJobs.RunAfter.LT_EQ(postgres.TimestampzT(time.Now())).
			AND(table.PgqJobs.RetryWaits.GT(postgres.Int(0))).
			AND(table.PgqJobs.InRun),
	).LIMIT(int64(pgq.workLimits)).Sql()

	rows, err := tx.Query(ctx, chooseTasksQuery, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot choose tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.PgqJobs
	for rows.Next() {
		var task model.PgqJobs
		inForRows := rows.Scan(
			&task.ID,
			&task.CreatedAt,
			&task.TaskName,
			&task.Payload,
			&task.RunAfter,
			&task.RetryWaits,
			&task.InRun,
		)
		if inForRows != nil {
			return false, fmt.Errorf("cannot scan task: %w", inForRows)
		}
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		handleFunc, ok := pgq.tasks[TaskKind(task.TaskName)]
		if !ok {
			logrus.Warnf("skip task %s because it does not exists: id=%d", task.TaskName, task.ID)
			continue
		}
		task.InRun = true
		err := handleFunc.Handle(ctx, TaskKind(task.TaskName), task.Payload)
		if err != nil {
			task.Error = lo.ToPtr(err.Error())
			task.RunAfter = time.Now().Add(pgq.taskRunAfterDuration)
			task.RetryWaits -= 1
			logrus.Errorf("cannot complite task id %d, err: %s", task.ID, err.Error())
		} else {
			task.InRun = false
		}
		err = pgq.updateTask(ctx, tx, task)
		if err != nil {
			logrus.Error(err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("cannot commit result: %w", err)
	}

	return true, nil
}

func (pgq *PgQueue) updateTask(ctx context.Context, tx pgx.Tx, task model.PgqJobs) error {
	stmt, args := table.PgqJobs.UPDATE(table.PgqJobs.AllColumns).
		MODEL(task).
		WHERE(table.PgqJobs.ID.EQ(postgres.Int(int64(task.ID)))).
		Sql()

	_, err := pgq.db.Exec(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("cannot update task: %w", err)
	}
	return nil
}

func (pgq *PgQueue) Schedule(ctx context.Context, taskName TaskKind, payload []byte) error {
	stmt, args := table.PgqJobs.INSERT(table.PgqJobs.AllColumns.Except(table.PgqJobs.ID)).
		MODEL(model.PgqJobs{
			CreatedAt:  time.Now(),
			TaskName:   string(taskName),
			Payload:    payload,
			RunAfter:   time.Now(),
			RetryWaits: int32(pgq.attemptsLeftDefault),
			Error:      nil,
			InRun:      true,
		}).Sql()

	_, err := pgq.db.Exec(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("cannot insert task %s to run\npayload: %s\nErr: %w", taskName, string(payload), err)
	}
	return nil
}
