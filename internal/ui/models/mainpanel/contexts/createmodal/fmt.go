package createmodal

import "strings"

const (
	mysqlDB      = "mysql"
	postgresqlDB = "postgres"
	sqlite       = ".db"
)

func formatEmodjiBasedOnDB(dsn, name string) string {
	if strings.HasPrefix(dsn, mysqlDB) {
		return name + " 🐬"
	}

	if strings.HasPrefix(dsn, postgresqlDB) {
		return name + " 🐘"
	}

	if strings.Contains(dsn, sqlite) {
		return name + " 🪶"
	}

	return name
}
