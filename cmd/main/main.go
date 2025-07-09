package main

import (
	"github.com/Turalchik/authentication-service/internal/auth_service"
	"github.com/Turalchik/authentication-service/internal/database"
	"github.com/Turalchik/authentication-service/internal/handlers"
	"github.com/Turalchik/authentication-service/internal/redisdb"
	"github.com/Turalchik/authentication-service/internal/repo"
	"github.com/Turalchik/authentication-service/internal/token_revocation_store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
)

// @title Authentication Service
// @version 1.0
// @description Authentication Service with JWT tokens

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg, err := GetConfigFromEnv()
	if err != nil {
		log.Fatalf("can't load variables from environment with error: %v", err)
	}

	dsn := database.NewPostgresDSN()
	db, err := database.NewDatabase(dsn, "pgx")
	if err != nil {
		log.Fatalf("Can't create database: %v", err)
	}
	redisClient, err := redisdb.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Can't create redis client: %v", err)
	}

	repository := repo.NewRepo(db)
	revocationStore := token_revocation_store.NewTokenRevocationStore(redisClient, "")
	authService := auth_service.NewAuthService(repository, revocationStore, cfg.TTLAccessToken, cfg.JWTSecretKey, cfg.WebhookURL)
	handler := handlers.NewHttpHandler(authService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
