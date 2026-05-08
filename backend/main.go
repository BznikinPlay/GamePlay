package main

import (
	"console-rental/database"
	"console-rental/handlers"
	"console-rental/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    for i := 0; i < 30; i++ {
        err := database.InitDB()
        if err == nil {
            log.Println("Successfully connected to database")
            break
        }
        log.Printf("Failed to connect to database (attempt %d/30): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }
    
    r := gin.Default()
    
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    r.POST("/api/register", handlers.Register)
    r.POST("/api/login", handlers.Login)
    r.GET("/api/consoles", handlers.GetConsoles)
    r.GET("/api/consoles/:id", handlers.GetConsoleByID)
    
    auth := r.Group("/api")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.GET("/user/profile", handlers.GetUserProfile)
        auth.PUT("/user/profile", handlers.UpdateUserProfile)
        
        auth.POST("/rentals", handlers.CreateRental)
        auth.GET("/my-rentals", handlers.GetUserRentals)
        auth.PUT("/rentals/:id/return", handlers.ReturnConsole)
    }
    
    log.Fatal(r.Run(":8080"))
}