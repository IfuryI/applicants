package main

import (
	http_check_cert "bitbucket.org/projectiu7/backend/src/master/internal/check_cert/delivery"
	repository_check_cert "bitbucket.org/projectiu7/backend/src/master/internal/check_cert/repository"
	service_check_cert "bitbucket.org/projectiu7/backend/src/master/internal/check_cert/service"
	http_dictionary "bitbucket.org/projectiu7/backend/src/master/internal/dictionary/delivery"
	repository_dictionary "bitbucket.org/projectiu7/backend/src/master/internal/dictionary/repository"
	service_dictionary "bitbucket.org/projectiu7/backend/src/master/internal/dictionary/service"
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/proto"
	"bitbucket.org/projectiu7/backend/src/master/internal/queue"
	sessionsDelivery "bitbucket.org/projectiu7/backend/src/master/internal/sessions/delivery"
	http_token_confirm "bitbucket.org/projectiu7/backend/src/master/internal/token_confirm/delivery"
	http_token_info "bitbucket.org/projectiu7/backend/src/master/internal/token_info/delivery"
	http_token_new "bitbucket.org/projectiu7/backend/src/master/internal/token_new/delivery"
	"bitbucket.org/projectiu7/backend/src/master/internal/users"
	usersHttp "bitbucket.org/projectiu7/backend/src/master/internal/users/delivery/http"
	usersDBStorage "bitbucket.org/projectiu7/backend/src/master/internal/users/repository/dbstorage"
	usersUseCase "bitbucket.org/projectiu7/backend/src/master/internal/users/usecase"
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Запуск сервера
func main() {
	app := NewApp()

	if err := app.Run(constants.Port); err != nil {
		log.Fatal(err)
	}
}

// App структура главного приложения
type App struct {
	server            *http.Server
	usersUC           users.UseCase
	serviceCheckCert  *service_check_cert.Service
	serviceDictionary *service_dictionary.Service
	serviceQueue      *queue.Service
	authMiddleware    middleware.Auth
	csrfMiddleware    middleware.Csrf
	logger            *logger.Logger
	sessionsDL        *sessionsDelivery.AuthClient
	sessionsConn      *grpc.ClientConn
	fileServer        proto.FileServerHandlerClient
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// NewApp инициализация приложения
func NewApp() *App {
	accessLogger := logger.NewAccessLogger()

	connStr := "postgres://mdb:mdb@localhost:5432/mdb"

	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	sessionsGrpcConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", constants.AuthPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to grpc auth server: %v\n", err)
	}
	sessionsDL := sessionsDelivery.NewAuthClient(sessionsGrpcConn)

	fileServerGrpcConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", constants.FileServerPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to grpc file server: %v\n", err)
	}
	fileServerService := proto.NewFileServerHandlerClient(fileServerGrpcConn)

	usersRepo := usersDBStorage.NewUserRepository(dbpool)
	checkCertRepo := repository_check_cert.NewCheckCertRepository(dbpool)
	dictionaryRepo := repository_dictionary.NewDictionaryRepository(dbpool)
	queueRepoProduce := queue.NewRepository(dbpool)

	usersUC := usersUseCase.NewUsersUseCase(usersRepo)
	checkCertService := service_check_cert.NewCheckCertUseCase(checkCertRepo)
	dictionaryService := service_dictionary.NewDictionaryService(dictionaryRepo)
	queueServiceProduce := queue.NewService(queueRepoProduce, nil)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, sessionsDL)
	csrfMiddleware := middleware.NewCsrfMiddleware(accessLogger)

	return &App{
		usersUC:           usersUC,
		serviceCheckCert:  checkCertService,
		serviceDictionary: dictionaryService,
		serviceQueue:      queueServiceProduce,
		authMiddleware:    authMiddleware,
		csrfMiddleware:    csrfMiddleware,
		logger:            accessLogger,
		sessionsDL:        sessionsDL,
		sessionsConn:      sessionsGrpcConn,
		fileServer:        fileServerService,
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Run запуск приложения
func (app *App) Run(port string) error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.Use(middleware.AccessLogMiddleware(app.logger))

	router.Static("/avatars", constants.AvatarsFileDir)

	router.Use(gin.Recovery())
	router.GET("/metrics", prometheusHandler())

	api := router.Group("/api")

	usersHttp.RegisterHTTPEndpoints(api, app.usersUC, app.sessionsDL, app.authMiddleware, app.fileServer, app.logger)
	http_check_cert.RegisterHTTPEndpoints(api, app.serviceCheckCert, app.authMiddleware, app.logger)
	http_dictionary.RegisterHTTPEndpoints(api, app.serviceDictionary, app.authMiddleware, app.logger)
	http_token_new.RegisterHTTPEndpoints(api, app.serviceQueue, app.authMiddleware, app.logger)
	http_token_info.RegisterHTTPEndpoints(api, app.serviceQueue, app.authMiddleware, app.logger)
	http_token_confirm.RegisterHTTPEndpoints(api, app.serviceQueue, app.authMiddleware, app.logger)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := app.server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to listen and serve: ", err)
		}
	}()

	// using graceful shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	_ = app.sessionsConn.Close()
	return app.server.Shutdown(ctx)
}
