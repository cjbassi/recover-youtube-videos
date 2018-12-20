package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"github.com/cjbassi/recover-youtube-videos/backend/src/api/youtube"
	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
	"github.com/cjbassi/recover-youtube-videos/backend/src/utils"
)

var (
	db    *database.DB
	dbURI string
)

func init() {
	godotenv.Load()
	dbURI = os.Getenv("DB_URI")

	db, err := database.Setup(dbURI)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
}

type body struct {
	AccessToken string `json:"access_token"`
}

func recoverVideos(event body) ([]youtube.Playlist, error) {
	service, err := youtube.Setup(event.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize youtube api service: %v", err)
	}
	playlists, err := service.FetchAllVideos()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all videos: %v", err)
	}
	playlistsOfRemovedVideos := func() []youtube.Playlist {
		var playlistsOfRemovedVideos []youtube.Playlist
		for i := 0; i < len(playlists); i++ {
			removed, notRemoved := utils.SplitPlaylistItems(playlists[i].PlaylistItems)
			if len(removed) > 0 {
				playlists[i].PlaylistItems = notRemoved
				newPlaylist := playlists[i]
				newPlaylist.PlaylistItems = removed
				playlistsOfRemovedVideos = append(playlistsOfRemovedVideos, newPlaylist)
			}
		}
		return playlistsOfRemovedVideos
	}()

	if db != nil {
		// check database for any video matches and replace the item in the playlist before we send it
		for i := 0; i < len(playlistsOfRemovedVideos); i++ {
			for j := 0; j < len(playlistsOfRemovedVideos[i].PlaylistItems); j++ {
				title := playlistsOfRemovedVideos[i].PlaylistItems[j].Title
				var video database.Video
				db.Where("id = ?", title).First(&video)
				if video.ID != "" {
					playlistsOfRemovedVideos[i].PlaylistItems[j].Title = video.Title
				}
			}
		}
		for _, playlist := range playlists {
			// if s.Database.Connection.NewRecord(playlist) {
			db.Create(&playlist)
			// }
			for _, video := range playlist.PlaylistItems {
				// if s.Database.Connection.NewRecord(video) {
				db.Create(&video)
				// }
			}
		}
	}

	return playlistsOfRemovedVideos, nil
}

func main() {
	lambda.Start(recoverVideos)
}
