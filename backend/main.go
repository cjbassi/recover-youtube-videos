package main

import (
	"os"
	"os/signal"
	"syscall"

	prefixed "github.com/cjbassi/logrus-prefixed-formatter"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
	"github.com/cjbassi/recover-youtube-videos/backend/src/server"
)

var (
	port            string
	frontendURL     string
	clientID        string
	databaseURL     string
	disableDatabase bool
)

func loadEnv() {
	env := os.Getenv("BACKEND_ENV")
	if env == "" {
		env = "development"
	}
	log.WithFields(log.Fields{
		"BACKEND_ENV": env,
	}).Info()

	godotenv.Load(".env." + env)
	godotenv.Load()

	port = ":" + os.Getenv("PORT")
	databaseURL = os.Getenv("DATABASE_URL")
	frontendURL = os.Getenv("FRONTEND_URL")
	clientID = os.Getenv("CLIENT_ID")
	disableDatabase = os.Getenv("DISABLE_DATABASE") == "true"
	log.WithFields(log.Fields{
		"DISABLE_DATABASE": disableDatabase,
	}).Info()
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&prefixed.TextFormatter{
		PrefixPadding:   9,
		TimestampFormat: "2006/01/02 15:04:05",
		FullTimestamp:   true,
	})
	databaseLogger := log.WithFields(log.Fields{
		"prefix": "database",
	})
	serverLogger := log.WithFields(log.Fields{
		"prefix": "server",
	})

	loadEnv()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var db *database.Database = &database.Database{}
	if !disableDatabase {
		db, err := database.Setup(databaseLogger, databaseURL)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		defer db.Close()
	} else {
		db = nil
	}

	s := server.Setup(serverLogger, port, db, clientID, frontendURL)

	go func() {
		<-quit
		s.Shutdown()
	}()

	s.ListenAndServe()
}
