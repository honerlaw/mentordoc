package server

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"server/controller"
	"server/service"
)

func main() {
	db := newDb()

	userDaoService := service.NewUserDao(db)
	userService := service.NewUser(userDaoService)
	userController := controller.NewUser(userService)

	router := chi.NewRouter()
	router.Route("/v1", func (r chi.Router) {
		userController.RegisterRoutes(r);
	})

	err := http.ListenAndServe(":5050", router)
	if err != nil {
		log.Fatal(err)
	}
}

func newDb() *sql.DB {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?autocommit=true",
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
	return db;
}
