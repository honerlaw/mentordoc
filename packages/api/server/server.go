package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
	"net/http"
	"os"
	"sync"
)

func StartServer(waitGroup *sync.WaitGroup) *http.Server {
	db := util.NewDb()

	// utilities
	transactionManager := util.NewTransactionManager(db, nil)
	aclService := acl.NewAclService(transactionManager, db, nil)
	authenticationService := NewAuthenticationService()
	validatorService := util.NewValidatorService()

	// repositories
	organizationRepository := NewOrganizationRepository(db, nil)
	userRepositoryService := NewUserRepository(db, nil)
	folderRepository := NewFolderRepository(db, nil)

	// services
	organizationService := NewOrganizationService(organizationRepository)
	userService := NewUserService(userRepositoryService, organizationService, transactionManager)
	folderService := NewFolderService(folderRepository, organizationService, aclService)

	// middlewares
	authenticationMiddleware := NewAuthenticationMiddleware(authenticationService, userService)

	// controllers
	userController := NewUserController(userService, validatorService, authenticationService)
	folderController := NewFolderController(validatorService, folderService, authenticationMiddleware)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/v1", func(r chi.Router) {
		userController.RegisterRoutes(r)
		folderController.RegisterRoutes(r)
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