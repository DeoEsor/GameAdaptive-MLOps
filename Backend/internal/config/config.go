package config

import (
	"os"
)

type configName string

const (
	DbDsn        = configName("db_dsn")
	ListenPort   = configName("listen_port")
	KafkaBrokers = configName("kafka_brokers")
)

func GetValue(cnfg configName) string {
	return os.Getenv(string(cnfg))
}
