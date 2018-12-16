package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			s.logger.Warnf("unauthorized request: no 'token' cookie")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return s.JWTKey, nil
		})
		if err != nil {
			s.logger.Errorf("failed to parse JWT: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ID := claims["ID"].(string)
			if ID == "" {
				s.logger.Warnf("unauthorized request: no 'ID' in JWT")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			r = r.WithContext(context.WithValue(context.Background(), "ID", ID))
			next.ServeHTTP(w, r)
		} else {
			s.logger.Warnf("unauthorized request: invalid JWT token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
	})
}
