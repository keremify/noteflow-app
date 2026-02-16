package main

import (
	"fmt"
	"log"

	"saasproject/internal/config"
	"saasproject/internal/database"
	"saasproject/internal/handlers"
	"saasproject/internal/models"
	"saasproject/internal/repository"
	"saasproject/internal/routes"
	"saasproject/internal/services"
)

func main() {
	if err := config.Set(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database.Connect()

	// Connect to database
	db := database.Connection()

	// Auto-migrate
	if err := db.AutoMigrate(&models.User{}, &models.Note{}, &models.Tag{}, &models.NoteTag{}, &models.RefreshToken{}); err != nil {
		log.Printf("warning: automigrate failed: %v", err)
	}

	// Repositories
	noteRepo := repository.NewNoteRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Services
	noteService := services.NewNoteService(noteRepo)
	userService := services.NewAuthService(userRepo, refreshTokenRepo)

	// Handlers
	noteHandler := handlers.NewNoteHandler(noteService)
	authHandler := handlers.NewAuthHandler(userService)

	// Routes
	r := routes.SetupRoutes(noteHandler, authHandler)

	// Server config
	cfg := config.Get()
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
