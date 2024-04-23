package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func getHashHandler(c *gin.Context) {
	// Retrieve the hash value from the environment variable
	hashValue := os.Getenv("HASH")
	if hashValue == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HASH environment variable not set"})
		return
	}

	// Return the hash value as JSON response
	c.JSON(http.StatusOK, gin.H{"hash": hashValue})
}

func healthCheckHandler(c *gin.Context) {
	// Respond with a simple message indicating the service is up
	c.String(http.StatusOK, "Service is up and running!")
}

func main() {
	// Initialize gin router
	router := gin.Default()

	// Define HTTP handlers for the endpoints
	router.GET("/hash", getHashHandler)
	router.GET("/health", healthCheckHandler)

	// Start the HTTP server
	fmt.Println("Server listening on :8081")
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
