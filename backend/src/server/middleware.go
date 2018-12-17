package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type middleware func(http.Handler) http.Handler
type middlewares []middleware

func (mws middlewares) apply(handler http.Handler) http.Handler {
	if len(mws) == 0 {
		return handler
	}
	return mws[1:].apply(mws[0](handler))
}

func (s *Server) logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			s.logger.WithFields(logrus.Fields{
				"requestID":   requestID,
				"method":      r.Method,
				"url":         r.URL.Path,
				"remoteaAddr": r.RemoteAddr,
				"userAgent":   r.UserAgent(),
				"duration":    time.Since(start),
			}).Info()
		}(time.Now())
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) tracing(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			id, err := s.sonyflake.NextID()
			if err != nil {
				s.logger.Errorf("failed to generate a new sonyflake id: %v", err)
				requestID = ""
			} else {
				requestID = strconv.Itoa(int(id))
			}
		}
		w.Header().Set("X-Request-Id", requestID)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) cors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.frontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		handler.ServeHTTP(w, r)
	})
}
