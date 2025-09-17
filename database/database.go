package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func NewDatabaseConnection() *sql.DB {
	db, connectErr := sql.Open("mysql", "demouser:demouserpassword@/project_horizon")
	if connectErr != nil {
		panic(connectErr)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	return db
}
