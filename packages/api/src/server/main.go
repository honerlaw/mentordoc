package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
	"net/http"
	"os"
	"server/controller"
	"server/service"
	"time"
)

func main() {
	db := newDb()

	validatorService := service.NewValidator()
	userDaoService := service.NewUserDao(db)
	userService := service.NewUser(userDaoService)
	userController := controller.NewUser(userService, validatorService)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/v1", func(r chi.Router) {
		userController.RegisterRoutes(r);
	})

	err := http.ListenAndServe(":5050", router)
	if err != nil {
		log.Fatal(err)
	}
}

func newDb() *sql.DB {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	log.Print(dataSource)

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
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", os.Getenv("DATABASE_NAME"), driver)
	if err != nil {
		log.Fatal(err)
	}

	err = migrator.Up()
	if err != nil {
		log.Fatal(err)
	}
}
