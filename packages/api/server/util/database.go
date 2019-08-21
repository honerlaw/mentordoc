package util

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"log"
	"os"
	"time"
)

func NewDb() *sql.DB {

	// multi statements is needed for migrations
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?autocommit=true&multiStatements=true",
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
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", os.Getenv("MIGRATION_DIR")), os.Getenv("DATABASE_NAME"), driver)

	if err != nil {
		log.Fatal(err)
	}

	err = migrator.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Print("finished running database migrations")
}
