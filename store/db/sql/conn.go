package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getConnection(database, endpoint string) (*gorm.DB, error) {
	var connector gorm.Dialector

	switch database {
	case "postgres":
		connector = postgres.Open(endpoint)
	case "mysql":
		connector = mysql.Open(endpoint)
	default:
		connector = sqlite.Open(endpoint)
	}

	return gorm.Open(connector, &gorm.Config{})
}
