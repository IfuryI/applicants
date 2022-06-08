package main

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/logger"
	"bitbucket.org/projectiu7/backend/src/master/internal/middleware"
	"bitbucket.org/projectiu7/backend/src/master/internal/modeling"
	constants "bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Запуск сервера
func main() {
	app := NewApp()

	if err := app.Run(constants.Port2); err != nil {
		log.Fatal(err)
	}
}

// App структура главного приложения
type App struct {
	server          *http.Server
	serviceModeller *modeling.Service
	logger          *logger.Logger
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// NewApp инициализация приложения
func NewApp() *App {
	accessLogger := logger.NewAccessLogger()
	broker := modeling.NewServer()

	modellerService := modeling.NewModeller(broker)

	return &App{
		serviceModeller: modellerService,
		logger:          accessLogger,
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

	api := router.Group("/api")

	modeling.RegisterHTTPEndpoints(api, app.serviceModeller, app.logger)

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

	return app.server.Shutdown(ctx)
}
