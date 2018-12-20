package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"github.com/cjbassi/recover-youtube-videos/backend/src/api/youtube"
	"github.com/cjbassi/recover-youtube-videos/backend/src/database"
)

var (
	db              *database.DB
	dbURI           string
	cors            = header{"Access-Control-Allow-Origin": "*"}
	jsonContentType = header{"Content-Type": "application/json"}
)

type body struct {
	AccessToken string `json:"accessToken"`
}

type request events.APIGatewayProxyRequest
type response events.APIGatewayProxyResponse
type header map[string]string
type playlistItems []youtube.PlaylistItem
type playlists []youtube.Playlist

func serverError() (response, error) {
	return response{
		StatusCode: http.StatusInternalServerError,
		Headers:    cors,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (response, error) {
	return response{
		StatusCode: status,
		Headers:    cors,
		Body:       http.StatusText(status),
	}, nil
}

func splitPlaylistItems(items playlistItems) (playlistItems, playlistItems) {
	removed := make(playlistItems, 0)
	notRemoved := make(playlistItems, 0)
	for _, item := range items {
		if item.Title == "Deleted video" || item.Title == "Private video" {
			removed = append(removed, item)
		} else {
			notRemoved = append(notRemoved, item)
		}
	}
	return removed, notRemoved
}

func combineHeaders(m1, m2 header) header {
	for key, val := range m1 {
		m2[key] = val
	}
	return m2
}

func fetchRemovedVideos(r request) (response, error) {
	var body body
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		log.Printf("failed to parse request body: %v", err)
		return clientError(http.StatusBadRequest)
	}

	service, err := youtube.Setup(body.AccessToken)
	if err != nil {
		log.Printf("failed to initialize youtube api service: %v", err)
		return serverError()
	}

	_playlists, err := service.FetchAllVideos()
	if err != nil {
		log.Printf("failed to fetch all videos: %v", err)
		return serverError()
	}

	playlistsOfRemovedVideos := func() playlists {
		var playlistsOfRemovedVideos playlists
		for i := 0; i < len(_playlists); i++ {
			removed, notRemoved := splitPlaylistItems(_playlists[i].PlaylistItems)
			if len(removed) > 0 {
				_playlists[i].PlaylistItems = notRemoved
				newPlaylist := _playlists[i]
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

		// add videos to database
		for _, playlist := range _playlists {
			for _, playlistItem := range playlist.PlaylistItems {
				video := database.Video{
					ID:        playlistItem.ID,
					Title:     playlistItem.Title,
					Thumbnail: playlistItem.Thumbnail,
				}
				db.Create(&video) // TODO
			}
		}
	}

	bytes, err := json.Marshal(playlistsOfRemovedVideos)
	if err != nil {
		log.Printf("failed to marshall removedVideos: %v", err)
		return serverError()
	}

	return response{
		StatusCode: http.StatusOK,
		Body:       string(bytes),
		Headers:    combineHeaders(jsonContentType, cors),
	}, nil
}

func main() {
	godotenv.Load()
	dbURI = os.Getenv("DB_URI")

	var err error
	db, err = database.Setup(dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
	} else {
		defer db.Close()
	}

	lambda.Start(fetchRemovedVideos)
}
