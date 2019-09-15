package util

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rubenv/sql-migrate"
	"log"
	"os"
	"time"
)

func NewDb() *sql.DB {

	// multi statements is needed for migrations
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?autocommit=true&multiStatements=true&parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Second)
	db.SetMaxIdleConns(0)

	runMigration(db)

	return db;
}

func runMigration(db *sql.DB) {
	// attempt to establish a connection before running the migrations
	for tick := 0; tick < 15; tick++ {
		err := db.Ping()
		if err == nil {
			break;
		}
		time.Sleep(time.Second)
	}

	log.Print("starting to run database migrations")

	migrations := &migrate.FileMigrationSource{
		Dir: os.Getenv("MIGRATION_DIR"),
	}

	_, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("finished running database migrations")
}
