package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/pokemon-chatbot-api/internal/config"
	"github.com/yourusername/pokemon-chatbot-api/internal/handlers"
	"github.com/yourusername/pokemon-chatbot-api/internal/repository"
	"github.com/yourusername/pokemon-chatbot-api/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize config
	cfg := config.New()

	// Initialize database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	pokemonService := services.NewPokemonService()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	pokemonHandler := handlers.NewPokemonHandler(pokemonService)

	// Setup router
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.GET("/:telegramId", userHandler.GetUser)
			users.GET("/:telegramId/check", userHandler.CheckRegistration)
		}

		// Pokemon routes
		pokemon := api.Group("/pokemon")
		{
			pokemon.GET("/:name", pokemonHandler.GetPokemon)
			pokemon.GET("/search/:query", pokemonHandler.SearchPokemon)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
