package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devforma/ultimate/cmd/api/handlers"
	"github.com/devforma/ultimate/internal/config"
	"github.com/devforma/ultimate/internal/database"
	"github.com/devforma/ultimate/internal/web"
	"github.com/streadway/amqp"
)

func main() {
	//init mysql db
	db, err := InitDB("config.yaml")
	if err != nil {
		log.Fatalf("init db failed: %v", err)
	}
	defer db.Close()

	//init rabbitmq connection
	mqConn, err := InitMQConn("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Fatalf("init mq connection failed: %v", err)
	}
	defer mqConn.Close()

	//init logger
	logger, err := InitLogger("api.log")
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}

	httpServer := InitAPIServer(db, mqConn, logger)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("stop server, new connections will not be accepted")
			} else {
				log.Fatalf("http server start failed: %v", err)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	fmt.Println("stoping server....")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown failed: %v\n", err)
		httpServer.Close()
	}
}

// InitLogger get logger
func InitLogger(logPath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return log.New(logFile, "[api]", log.LstdFlags), nil
}

// InitMQConn get mq connection
func InitMQConn(mqUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// InitDB get database connection
func InitDB(cfgPath string) (*database.DB, error) {
	var cfg database.Config
	if err := config.Parse(cfgPath, &cfg); err != nil {
		return nil, err
	}

	db, err := database.Open(&cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitAPIServer init http.Server
func InitAPIServer(db *database.DB, mqConn *amqp.Connection, logger *log.Logger) http.Server {
	handler := web.NewHandler("404 not found HaHa")

	canalAPI := handlers.CanalAPI{DB: db}

	handler.AddRoute("GET", "/canal", canalAPI.SingleCanal)
	handler.AddRoute("GET", "/canals", canalAPI.ListCanal)

	mqAPI := handlers.MQAPI{Logger: logger, Conn: mqConn}
	handler.AddRoute("GET", "/mq/workqueue", mqAPI.WorkQueue)

	return http.Server{
		Addr:         ":80",
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
}
