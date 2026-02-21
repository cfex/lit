package connections

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB(driver, dsn string) {
	var err error
	DB, err = sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
}
func CleanupDB() {
	DB.Close()
}
