package app

import (
	"fmt"
	"forum/internal/controller/http1"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"forum/pkg/config"
	"forum/pkg/database"
	"io"
	"log"
	"os"
)

const secret string = "secret"

func Run(cfg *config.Conf) {
	// Prepare logger
	file, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("cannot create log file: %v", err)
	}
	defer file.Close()
	logWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(logWriter)

	// Prepare database
	db, err := database.ConnectSqlte(&cfg.Database)
	if err != nil {
		log.Fatalf("error occured while connecting database: %s", err.Error())
		return
	}
	// Close connection database
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal("can't close connection db, err:", err)
		} else {
			log.Println("db closed")
		}
	}()

	// Prepare router <- -> service  <- -> repository
	repo := repository.NewRepository(db)
	service := service.NewService(repo, secret)
	handler := http1.NewHandler(service, secret)
	server := new(server.Server)
	// Start listening server
	log.Fatalf("error occured while listening server: %s", server.Run(&cfg.API, handler.InitRoutes(cfg)))
}
