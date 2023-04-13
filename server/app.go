package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"health/routes"
	service_email "health/services/email"
)

type App struct {
	httpServer *http.Server

	db     *mongo.Database
	mailer *service_email.Mailer
}

type EmailContent struct {
	Name string
}

func InitApp() *App {
	db := initDB()

	mailer := service_email.NewMailer()

	return &App{
		db:     db,
		mailer: mailer,
	}
}

func (app *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	routes.InitRoutes(router, app.db, app.mailer)

	// Конфиги для сервера
	app.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("App started, listen port: %s", port)

		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	// Ждем когда к нам придет сигнал
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer log.Println("App shutdown ...")
	defer shutdown()

	return app.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	dbURI := viper.GetString("db.uri")

	// креды для авторизации
	credential := options.Credential{
		Username: viper.GetString("db.options.user"),
		Password: viper.GetString("db.options.password"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI).SetAuth(credential))
	if err != nil {
		log.Fatalf("Error occurred while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("App connected to db_uri %s", dbURI)

	return client.Database(viper.GetString("db.name"))
}