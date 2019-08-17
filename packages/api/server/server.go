package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/honerlaw/mentordoc/server/transaction"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func StartServer(waitGroup *sync.WaitGroup) *http.Server {
	db := newDb()

	validatorService := NewValidatorService()
	transactionManager := transaction.NewTransactionManager(db)
	organizationRepository := NewOrganizationRepository(db, nil)
	organizationService := NewOrganizationService(organizationRepository)
	userRepositoryService := NewUserRepository(db, nil)
	userService := NewUserService(userRepositoryService, organizationService, transactionManager)
	userController := NewUserController(userService, validatorService)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/v1", func(r chi.Router) {
		userController.RegisterRoutes(r);
	})

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: router,
	}

	go func() {
		if waitGroup != nil {
			defer waitGroup.Done()
		}
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	log.Print("successfully started server")

	return server
}

func StopServer(server *http.Server) {
	err := server.Shutdown(context.Background());
	if err != nil {
		panic(err)
	}
}

func newDb() *sql.DB {

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

