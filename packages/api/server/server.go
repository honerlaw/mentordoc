package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/util"
	"log"
	"net/http"
	"os"
	"sync"
)

type Server struct {
	HttpServer                *http.Server
	TransactionManager        *util.TransactionManager
	AclService                *acl.AclService
	AuthenticationService     *AuthenticationService
	ValidatorService          *util.ValidatorService
	OrganizationRepository    *OrganizationRepository
	UserRepository            *UserRepository
	FolderRepository          *FolderRepository
	DocumentRepository        *DocumentRepository
	DocumentContentRepository *DocumentContentRepository
	OrganizationService       *OrganizationService
	UserService               *UserService
	FolderService             *FolderService
	DocumentService           *DocumentService
	AuthenticationMiddleware  *AuthenticationMiddleware
	UserController            *UserController
	FolderController          *FolderController
	DocumentController        *DocumentController
}

func StartServer(waitGroup *sync.WaitGroup) *Server {
	db := util.NewDb()

	// utilities
	transactionManager := util.NewTransactionManager(db, nil)
	aclService := acl.NewAclService(transactionManager, db, nil)
	authenticationService := NewAuthenticationService()
	validatorService := util.NewValidatorService()

	// repositories
	organizationRepository := NewOrganizationRepository(db, nil)
	userRepository := NewUserRepository(db, nil)
	folderRepository := NewFolderRepository(db, nil)
	documentRepository := NewDocumentRepository(db, nil)
	documentContentRepository := NewDocumentContentRepository(db, nil)

	// services
	organizationService := NewOrganizationService(organizationRepository)
	userService := NewUserService(userRepository, organizationService, transactionManager)
	folderService := NewFolderService(folderRepository, organizationService, aclService)
	documentService := NewDocumentService(documentRepository, documentContentRepository, organizationService, folderService, aclService, transactionManager)

	// middlewares
	authenticationMiddleware := NewAuthenticationMiddleware(authenticationService, userService)

	// controllers
	userController := NewUserController(userService, validatorService, authenticationService, authenticationMiddleware)
	folderController := NewFolderController(validatorService, folderService, authenticationMiddleware, aclService)
	documentController := NewDocumentController(validatorService, documentService, authenticationMiddleware, aclService)

	err := aclService.Init()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/v1", func(r chi.Router) {
		userController.RegisterRoutes(r)
		folderController.RegisterRoutes(r)
		documentController.RegisterRoutes(r)
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
		AuthenticationService:     authenticationService,
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
	}
}

func StopServer(server *Server) {
	err := server.HttpServer.Shutdown(context.Background());
	if err != nil {
		panic(err)
	}
}
