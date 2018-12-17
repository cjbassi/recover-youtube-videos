package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

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

func (s *Server) fetchRemovedVideos(w http.ResponseWriter, r *http.Request) {
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
