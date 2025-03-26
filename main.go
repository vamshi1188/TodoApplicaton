package main

import (
	"log"
)

func main() {
	// Initialize AWS S3 client
	initS3Client()

	// Load tasks from S3
	loadTasksFromS3()

	// Initialize and run the UI
	if err := initializeUI(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
