package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/futurenda/google-auth-id-token-verifier"

	"github.com/cjbassi/recover-youtube-videos/backend/src/api"
	. "github.com/cjbassi/recover-youtube-videos/backend/src/models"
	. "github.com/cjbassi/recover-youtube-videos/backend/src/utils"
)

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	if h := atomic.LoadInt64(&s.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.Write([]byte(fmt.Sprintf("uptime: %s", time.Since(time.Unix(0, h)))))
	}
}

func (s *Server) fetchMissingVideos(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		AccessToken string `json:"access_token"`
	}
	err := decoder.Decode(&body)
	if err != nil {
		s.logger.Errorf("failed to parse request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	service, err := api.Setup(body.AccessToken)
	if err != nil {
		s.logger.Errorf("failed to initialize youtube api service: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playlists, err := service.FetchAllVideos()
	if err != nil {
		s.logger.Errorf("failed to fetch all videos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playlistsOfRemovedVideos := func() []Playlist {
		var playlistsOfRemovedVideos []Playlist
		for i := 0; i < len(playlists); i++ {
			removed, notRemoved := SplitVideos(playlists[i].Videos)
			if len(removed) > 0 {
				playlists[i].Videos = notRemoved
				newPlaylist := playlists[i]
				newPlaylist.Videos = removed
				playlistsOfRemovedVideos = append(playlistsOfRemovedVideos, newPlaylist)
			}
		}
		return playlistsOfRemovedVideos
	}()
	if s.Database != nil {
		// check database for any video matches and replace the item in the playlist before we send it
		for i := 0; i < len(playlistsOfRemovedVideos); i++ {
			for j := 0; j < len(playlistsOfRemovedVideos[i].Videos); j++ {
				title := playlistsOfRemovedVideos[i].Videos[j].Title
				var video Video
				s.Database.Connection.Where("id = ?", title).First(&video)
				if video.ID != "" {
					playlistsOfRemovedVideos[i].Videos[j] = video
				}
			}
		}
		for _, playlist := range playlists {
			if s.Database.Connection.NewRecord(playlist) {
				s.Database.Connection.Create(&playlist)
			}
			for _, video := range playlist.Videos {
				if s.Database.Connection.NewRecord(video) {
					s.Database.Connection.Create(&video)
				}
			}
		}
	}
	json.NewEncoder(w).Encode(
		playlistsOfRemovedVideos,
	)
}

func (s *Server) tokenSignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		IDToken string `json:"idtoken"`
	}
	err := decoder.Decode(&body)
	if err != nil {
		s.logger.Errorf("failed to parse request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := googleAuthIDTokenVerifier.Verifier{}
	err = v.VerifyIDToken(body.IDToken, []string{
		s.clientID,
	})
	if err != nil {
		s.logger.Errorf("failed to validate token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	claimSet, err := googleAuthIDTokenVerifier.Decode(body.IDToken)
	if err != nil {
		s.logger.Errorf("failed to decode token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := User{
		Name: claimSet.Name,
		ID:   claimSet.Sub,
	}

	if s.Database != nil {
		if s.Database.Connection.NewRecord(user) {
			s.Database.Connection.Create(&user)
		}
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	expiration := time.Now().Add(s.tokenExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  user.ID,
		"exp": expiration,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(s.JWTKey)
	if err != nil {
		s.logger.Errorf("failed to sign token: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "token", Value: tokenString, Expires: expiration}
	http.SetCookie(w, &cookie)
}
