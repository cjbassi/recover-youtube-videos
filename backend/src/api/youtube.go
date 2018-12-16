package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"

	. "github.com/cjbassi/recover-youtube-videos/backend/src/models"
)

type YTService struct {
	service *youtube.Service
}

func Setup(access_token string) (*YTService, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	token := oauth2.Token{
		AccessToken: access_token,
	}
	client := config.Client(ctx, &token)
	service, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to create a youtube service from the client/token: %v", err)
	}
	return &YTService{service}, nil
}

func (s *YTService) FetchPlaylists() ([]Playlist, error) {
	playlists := []Playlist{}
	nextPageToken := ""
	for {
		playlistCall := s.service.Playlists.List("snippet").Mine(true).MaxResults(50).PageToken(nextPageToken)
		playlistResponse, err := playlistCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to call api: %v", err)
		}
		for _, playlist := range playlistResponse.Items {
			myPlaylist := Playlist{
				Title: playlist.Snippet.Title,
				ID:    playlist.Id,
				Channel: &Channel{
					ID:   playlist.Snippet.ChannelId,
					Name: playlist.Snippet.ChannelTitle,
				},
			}
			playlists = append(playlists, myPlaylist)
		}

		nextPageToken = playlistResponse.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return playlists, nil
}

func (s *YTService) FetchAllVideos() ([]Playlist, error) {
	playlists, err := s.FetchPlaylists()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlists: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(playlists))
	for i := 0; i < len(playlists); i++ {
		go func(i int) {
			_ = s.FetchPlaylistItems(&playlists[i]) // TODO
			wg.Done()
		}(i)
	}
	wg.Wait()
	return playlists, err
}

func (s *YTService) ChannelID() (string, error) {
	channelCall := s.service.Channels.List("id").Mine(true)
	channelResponse, err := channelCall.Do()
	if err != nil {
		return "", fmt.Errorf("failed to call api: %v", err)
	}
	return channelResponse.Items[0].Id, nil
}

func (s *YTService) FetchPlaylistItems(playlist *Playlist) error {
	nextPageToken := ""
	for {
		playlistItemsCall := s.service.PlaylistItems.List("snippet,contentDetails").PlaylistId(playlist.ID).MaxResults(50).PageToken(nextPageToken)
		playlistItemsResponse, err := playlistItemsCall.Do()
		if err != nil {
			return fmt.Errorf("failed to call api: %v", err)
		}
		for _, playlistItem := range playlistItemsResponse.Items {
			myPlaylistItem := Video{
				Title: playlistItem.Snippet.Title,
				ID:    playlistItem.ContentDetails.VideoId,
			}
			playlist.Videos = append(playlist.Videos, myPlaylistItem)
		}

		nextPageToken = playlistItemsResponse.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return nil
}
