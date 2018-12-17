package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"

	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
)

type Server struct {
	HTTPServer      http.Server
	port            string
	logger          *logrus.Entry
	sonyflake       *sonyflake.Sonyflake
	healthy         int64
	Database        *database.Database
	JWTKey          []byte
	clientID        string
	shutdownTime    time.Duration
	frontendURL     string
	tokenExpiration time.Duration
}

func Setup(logger *logrus.Entry, port string, db *database.Database, clientID string, frontendURL string) *Server {
	s := &Server{
		port:            port,
		logger:          logger,
		sonyflake:       sonyflake.NewSonyflake(sonyflake.Settings{}),
		healthy:         0,
		Database:        db,
		JWTKey:          []byte(fmt.Sprintf("%d", rand.Int63())),
		clientID:        clientID,
		shutdownTime:    15 * time.Second,
		frontendURL:     frontendURL,
		tokenExpiration: time.Hour * 24 * 14,
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthz", s.healthz)
	router.Handle("/fetchremovedvideos", http.HandlerFunc(s.fetchRemovedVideos)).Methods("POST")

	s.HTTPServer = http.Server{
		Addr:         port,
		Handler:      (middlewares{s.tracing, s.logging, s.cors}).apply(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Infof("Server is ready to handle requests at %q", port)
	atomic.StoreInt64(&s.healthy, time.Now().UnixNano())

	return s
}

func (s *Server) ListenAndServe() {
	s.logger.Infof("Server is starting...")
	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatalf("Could not listen on %q: %s", s.port, err)
	}
}

func (s *Server) Shutdown() {
	atomic.StoreInt64(&s.healthy, 0)
	s.logger.Infof("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTime)
	defer cancel()

	s.HTTPServer.SetKeepAlivesEnabled(false)
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("Could not gracefully shutdown the server: %s", err)
	}

	s.logger.Infof("Server stopped")
}
