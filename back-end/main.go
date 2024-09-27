package main

import (
	"your-backend-module/config"
	"your-backend-module/handlers"
	"your-backend-module/middlewares"
"time"
	"context" // Eksik olan context paketi eklendi
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},  // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup database connection
	client, userCollection, documentCollection := config.SetupDatabase()
	defer client.Disconnect(context.Background())

	// Routes
	r.GET("/hello-backend", handlers.HelloHandler)
	r.POST("/register", handlers.RegisterHandler(userCollection))
	r.POST("/login", handlers.LoginHandler(userCollection))
	r.POST("/upload", handlers.UploadHandler(documentCollection)) // Upload rotasına documentCollection eklendi
	r.GET("/validate-token", middleware.ValidateToken())          // Token doğrulama rotası
	r.GET("/documents", handlers.GetDocumentsHandler(documentCollection))
	r.DELETE("/delete-file/:id", handlers.DeleteFileHandler(documentCollection)) // Dosya silme rotası

	// Start the server
	r.Run(":8080")
}
