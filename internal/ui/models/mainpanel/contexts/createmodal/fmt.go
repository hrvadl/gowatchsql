package createmodal

import "strings"

const (
	mysqlDB      = "mysql"
	postgresqlDB = "postgres"
)

func formatEmodjiBasedOnDB(dsn, name string) string {
	if strings.HasPrefix(dsn, mysqlDB) {
		return name + " 🐬"
	}

	if strings.HasPrefix(dsn, postgresqlDB) {
		return name + " 🐘"
	}

	return name
}
