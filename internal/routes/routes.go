package routes

import (
	"saasproject/internal/handlers"
	"saasproject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	noteHandler *handlers.NoteHandler,
	authHandler *handlers.AuthHandler,
) *gin.Engine {

	r := gin.Default()

	// Frontend static files
	r.StaticFile("/", "./src/index.html")
	r.StaticFile("/index.html", "./src/index.html")
	r.StaticFile("/style.css", "./src/style.css")
	r.StaticFile("/script.js", "./src/script.js")

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.Refresh)
	r.POST("/logout", authHandler.Logout)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	userOrAdmin := middleware.RequireRole("user", "admin")

	api.POST("/notes", userOrAdmin, noteHandler.Create)
	api.GET("/notes", userOrAdmin, noteHandler.List)
	api.PUT("/notes/:id", userOrAdmin, noteHandler.Update)
	api.DELETE("/notes/:id", userOrAdmin, noteHandler.Delete)

	return r
}
