package main

import (
	"log"
	"net/http"
	"s3-interaction/config"
	"s3-interaction/router"
	"s3-interaction/s3Client"
)

func main() {

	// Load the appConfig
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading appConfig: %v", err)
	}

	// Initialize the S3 client
	s3, err := s3Client.NewS3Client(
		appConfig.AWS.Key, appConfig.AWS.Secret, appConfig.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Initialize the router
	r, err := router.NewRouter(s3)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	r.SetupRoutes()

	// Start the server
	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
