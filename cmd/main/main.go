package main

import (
	"github.com/Turalchik/authentication-service/internal/auth_service"
	"github.com/Turalchik/authentication-service/internal/database"
	"github.com/Turalchik/authentication-service/internal/handlers"
	"github.com/Turalchik/authentication-service/internal/repo"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	godotenv.Load()

	cfg, err := GetConfigFromEnv()
	if err != nil {
		log.Fatalf("can't load variables from environment with error: %v", err)
	}

	dsn := database.NewPostgresDSN()
	db, err := database.NewDatabase(dsn, "pgx")
	if err != nil {
		log.Fatalf("Can't create database: %v", err)
	}

	repository := repo.NewRepo(db)
	authService := auth_service.NewAuthService(repository, cfg.TTLAccessToken, cfg.JWTSecretKey, cfg.WebhookURL)
	handler := handlers.NewHttpHandler(authService)

	// 3. Настройка сервера
	server := &http.Server{
		Addr:    ":8080", // Порт
		Handler: handler,
	}

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
