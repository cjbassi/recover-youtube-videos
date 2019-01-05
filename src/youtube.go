package src

import (
	"fmt"
	"sync"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title string `json:"title"`
	ID    string `json:"id"`
}

type PlaylistItem struct {
	Title    string `json:"title"`
	ID       string `json:"id"`
	Position int64  `json:"position"`
}

type Playlist struct {
	Title         string         `json:"title"`
	ID            string         `json:"id"`
	PlaylistItems []PlaylistItem `json:"playlistItems"`
}

func fetchPlaylists(service *youtube.Service) ([]Playlist, error) {
	playlists := []Playlist{}
	nextPageToken := ""
	for {
		playlistResponse, err := service.Playlists.
			List("snippet").
			Mine(true).
			MaxResults(50).
			PageToken(nextPageToken).
			Do()
		if err != nil {
			return nil, fmt.Errorf("failed to execute api call: %v", err)
		}
		for _, item := range playlistResponse.Items {
			myPlaylist := Playlist{
				ID:    item.Id,
				Title: item.Snippet.Title,
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

func FetchAllPlaylistItems(service *youtube.Service) ([]Playlist, error) {
	playlists, err := fetchPlaylists(service)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlists: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(len(playlists))
	for i := 0; i < len(playlists); i++ {
		go func(i int) {
			_ = fetchPlaylistItems(service, &playlists[i]) // TODO
			wg.Done()
		}(i)
	}
	wg.Wait()
	return playlists, err
}

func fetchPlaylistItems(service *youtube.Service, playlist *Playlist) error {
	nextPageToken := ""
	for {
		playlistItemsResponse, err := service.PlaylistItems.
			List("snippet,contentDetails").
			PlaylistId(playlist.ID).
			MaxResults(50).
			PageToken(nextPageToken).
			Do()
		if err != nil {
			return fmt.Errorf("failed to execute api call: %v", err)
		}
		for _, item := range playlistItemsResponse.Items {
			myPlaylistItem := PlaylistItem{
				ID:       item.ContentDetails.VideoId,
				Title:    item.Snippet.Title,
				Position: item.Snippet.Position,
			}
			playlist.PlaylistItems = append(playlist.PlaylistItems, myPlaylistItem)
		}
		nextPageToken = playlistItemsResponse.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return nil
}
