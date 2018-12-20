package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
)

func hardMigrate() {
	godotenv.Load()
	dbURI := os.Getenv("DB_URI")

	db, err := database.Setup(dbURI)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	db.HardMigrate()
}

func main() {
	lambda.Start(hardMigrate)
}
