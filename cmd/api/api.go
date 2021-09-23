package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devforma/ultimate/internal/config"
	"github.com/devforma/ultimate/internal/database"
	"github.com/devforma/ultimate/internal/web"
)

func main() {
	db, err := InitDB("config.yaml")
	if err != nil {
		log.Fatalf("init db failed: %v", err)
	}
	defer db.Close()

	httpServer := InitAPIServer(db)

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

type Canal struct {
	ID       int    `db:"id" json:"id"`
	Title    string `db:"title" json:"title"`
	Duration int    `db:"duration" json:"duration"`
}

// InitAPIServer init http.Server
func InitAPIServer(db *database.DB) http.Server {
	handler := web.NewHandler("404 not found HaHa")

	handler.AddRoute("GET", "/canal", func(w http.ResponseWriter, r *http.Request) {
		var canal Canal
		db.Get(&canal, "SELECT * FROM `canal` WHERE `id`=?", r.URL.Query().Get("id"))

		time.Sleep(8 * time.Second)

		if data, err := json.Marshal(canal); err == nil {
			w.Write(data)
		}
	})

	handler.AddRoute("GET", "/canals", func(w http.ResponseWriter, r *http.Request) {
		var canals []Canal
		db.Select(&canals, "SELECT * FROM `canal`")
		if data, err := json.Marshal(canals); err == nil {
			w.Write(data)
		}
	})

	return http.Server{
		Addr:         ":80",
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
}
