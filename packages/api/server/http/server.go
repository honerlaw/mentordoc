package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/honerlaw/mentordoc/server/http/controller"
	middleware2 "github.com/honerlaw/mentordoc/server/http/middleware"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/document"
	"github.com/honerlaw/mentordoc/server/lib/folder"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/user"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
	"net/http"
	"os"
	"sync"
)

type Server struct {
	HttpServer                *http.Server
	TransactionManager        *util.TransactionManager
	AclService                *acl.AclService
	TokenService              *util.TokenService
	ValidatorService          *util.ValidatorService
	OrganizationRepository    *organization.OrganizationRepository
	UserRepository            *user.UserRepository
	FolderRepository          *folder.FolderRepository
	DocumentRepository        *document.DocumentRepository
	DocumentContentRepository *document.DocumentContentRepository
	OrganizationService       *organization.OrganizationService
	UserService               *user.UserService
	FolderService             *folder.FolderService
	DocumentService           *document.DocumentService
	AuthenticationMiddleware  *middleware2.AuthenticationMiddleware
	UserController            *controller.UserController
	FolderController          *controller.FolderController
	DocumentController        *controller.DocumentController
	OrganizationController    *controller.OrganizationController
}

func StartServer(waitGroup *sync.WaitGroup) *Server {
	db := util.NewDb()

	// utilities
	transactionManager := util.NewTransactionManager(db, nil)
	aclService := acl.NewAclService(transactionManager, db, nil)
	tokenService := util.NewTokenService()
	validatorService := util.NewValidatorService()

	// repositories
	organizationRepository := organization.NewOrganizationRepository(db, nil)
	userRepository := user.NewUserRepository(db, nil)
	folderRepository := folder.NewFolderRepository(db, nil)
	documentRepository := document.NewDocumentRepository(db, nil)
	documentContentRepository := document.NewDocumentContentRepository(db, nil)

	// services
	organizationService := organization.NewOrganizationService(organizationRepository, aclService)
	userService := user.NewUserService(userRepository, organizationService, transactionManager, aclService)
	folderService := folder.NewFolderService(folderRepository, organizationService, aclService)
	documentService := document.NewDocumentService(documentRepository, documentContentRepository, organizationService, folderService, aclService, transactionManager)

	// middlewares
	authenticationMiddleware := middleware2.NewAuthenticationMiddleware(tokenService, userService)

	// controllers
	userController := controller.NewUserController(userService, validatorService, tokenService, authenticationMiddleware)
	folderController := controller.NewFolderController(validatorService, folderService, authenticationMiddleware, aclService)
	documentController := controller.NewDocumentController(validatorService, documentService, authenticationMiddleware, aclService)
	organizationController := controller.NewOrganizationController(organizationService, authenticationMiddleware, aclService)

	err := aclService.Init()
	if err != nil {
		log.Fatal(err)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router := chi.NewRouter()
	router.Use(corsMiddleware.Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/v1", func(r chi.Router) {
		userController.RegisterRoutes(r)
		folderController.RegisterRoutes(r)
		documentController.RegisterRoutes(r)
		organizationController.RegisterRoutes(r)
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT")),
		Handler: router,
	}

	go func() {
		if waitGroup != nil {
			defer waitGroup.Done()
		}
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	log.Print("successfully started server")

	return &Server{
		HttpServer:                httpServer,
		TransactionManager:        transactionManager,
		AclService:                aclService,
		TokenService:              tokenService,
		ValidatorService:          validatorService,
		OrganizationRepository:    organizationRepository,
		UserRepository:            userRepository,
		FolderRepository:          folderRepository,
		DocumentRepository:        documentRepository,
		DocumentContentRepository: documentContentRepository,
		OrganizationService:       organizationService,
		UserService:               userService,
		FolderService:             folderService,
		DocumentService:           documentService,
		AuthenticationMiddleware:  authenticationMiddleware,
		UserController:            userController,
		FolderController:          folderController,
		DocumentController:        documentController,
		OrganizationController:    organizationController,
	}
}

func StopServer(server *Server) {
	err := server.HttpServer.Shutdown(context.Background());
	if err != nil {
		panic(err)
	}
}
