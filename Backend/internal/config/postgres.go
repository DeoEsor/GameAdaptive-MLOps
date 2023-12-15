package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPostgres подключение к базе данных. Возвращает указатель на коннект к базе и ошибку.
// Перед использованием обязательно нужно сделать парсинг env переменных через вызов функции flags.InitServiceFlags
func ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, GetValue(DbDsn))
}
