package main

import (
	"context"
	"go-api/db"
	"go-api/handlers"
	"go-api/middleware"
	"go-api/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to databases
	db.ConnectPG()
	db.ConnectMongo()

	// Ensure all tables are created if they don't exist yet
	if err := db.MigratePG(models.Customer{}, models.Session{}); err != nil {
		log.Fatalf("Failed to migrate postgres: %v", err)
	}
	log.Println("POSTGRES successfully initiated")

	if err := db.MigrateMongo("customers", "sessions"); err != nil {
		log.Fatalf("Failed to migrate mongo: %v", err)
	}
	log.Println("MONGODB successfully initiated")

	// Initiate handlers
	customerHander := handlers.NewCustomerHandler()
	authHandler := handlers.NewSessionHandler()

	// Initialise router
	if os.Getenv("GIN_MODE") != "production" {
		gin.SetMode(gin.DebugMode)
	}
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Ok!"})
	})

	// Initialize routes
	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		api.Use(middleware.JWTMiddleware())
		api.GET("/profile", authHandler.GetProfile)
		api.POST("/logout", authHandler.Logout)
		api.GET("/customers", customerHander.Get)
		api.GET("/customers/:id", customerHander.GetOne)
		api.GET("/customers/email/:email", customerHander.GetByEmail)
		api.PUT("/customers/:id", customerHander.Update)
		api.DELETE("/customers/:id", customerHander.SoftDelete)
		api.PATCH("/customers/:id", customerHander.HardDelete)
	}

	// Public homepage
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Serve static files
	router.Static("/static", "./static")

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("App stopped successfully")
}
