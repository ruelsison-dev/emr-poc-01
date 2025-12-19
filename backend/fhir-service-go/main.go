package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// FHIR basic patient endpoints (scaffold)
	patientRoutes := r.Group("/Patient")
	{
		patientRoutes.POST("/", createPatient)
		patientRoutes.GET("/:id", getPatient)
	}

	log.Printf("Starting FHIR service on :%s", port)
	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
