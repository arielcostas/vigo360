package database

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/logger"
)

var db *sqlx.DB

func start() {
	logger := logger.NewLogger("BBDD")
	var dsn string = os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_BASE")
	var err error
	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Critical("error connecting to mysql: %s", err.Error())
	}

	logger.Information("database connection established")

	err = conn.Ping()
	if err != nil {
		logger.Critical("couldn't ping database: %s", err.Error())
	}

	_, err = conn.Exec("SET lc_time_names = 'es_ES';")
	if err != nil {
		logger.Critical("error configuring database locale: %s", err.Error())
	}
	_, err = conn.Exec("SET @@session.time_zone='+00:00';")
	if err != nil {
		logger.Critical("error configuring database timezone: %s", err.Error())
	}
	logger.Information("database configured")
	db = conn
}

func GetDB() *sqlx.DB {
	if db == nil {
		start()
	}
	return db
}
