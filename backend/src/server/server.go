// https://gist.github.com/creack/4c00ee404f2d7bd5983382cc93af5147
// https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"

	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
)

type Server struct {
	HTTPServer      http.Server
	port            string
	logger          *logrus.Entry
	Database        *database.Database
	clientID        string
	shutdownTimeout time.Duration
	frontendURL     string
	stats           *stats.Stats
}

func Setup(logger *logrus.Entry, port string, db *database.Database, clientID string, frontendURL string) *Server {
	s := &Server{
		port:            port,
		logger:          logger,
		Database:        db,
		clientID:        clientID,
		shutdownTimeout: 15 * time.Second,
		frontendURL:     frontendURL,
		stats:           stats.New(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/stats", s.statsRoute)
	router.HandleFunc("/fetchremovedvideos", s.fetchRemovedVideosRoute).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowCredentials: true,
	})

	n := negroni.New().With(negroni.NewRecovery(), negroni.NewLogger(), c, s.stats)
	n.UseHandler(router)

	s.HTTPServer = http.Server{
		Addr:         port,
		Handler:      n,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Infof("Server is ready to handle requests at %q", port)

	return s
}

func (s *Server) ListenAndServe() {
	s.logger.Infof("Server is starting...")
	if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatalf("Could not listen on %q: %s", s.port, err)
	}
}

func (s *Server) Shutdown() {
	s.logger.Infof("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	s.HTTPServer.SetKeepAlivesEnabled(false)
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("Could not gracefully shutdown the server: %s", err)
	}

	s.logger.Infof("Server stopped")
}
