package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/youtube/v3"

	. "github.com/cjbassi/recover-youtube-videos/src"
)

var (
	LIBRARY_FILE            = filepath.Join(os.Args[1], "library.json")
	RECOVERED_VIDEOS_FILE   = filepath.Join(os.Args[1], "recovered_videos.json")
	UNRECOVERED_VIDEOS_FILE = filepath.Join(os.Args[1], "unrecovered_videos.json")
)

func partitionPlaylistItems(playlistItems []PlaylistItem) ([]PlaylistItem, []PlaylistItem) {
	removed := []PlaylistItem{}
	notRemoved := []PlaylistItem{}
	for _, playlistItem := range playlistItems {
		if playlistItem.Title == "Deleted video" || playlistItem.Title == "Private video" {
			removed = append(removed, playlistItem)
		} else {
			notRemoved = append(notRemoved, playlistItem)
		}
	}
	return removed, notRemoved
}

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	token := GetToken(config)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	client := config.Client(context.Background(), token)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	libraryBytes, err := ioutil.ReadFile(LIBRARY_FILE)
	if err != nil {
		log.Fatalf("failed to read library file: %v", err)
	}
	var library []Video
	err = json.Unmarshal(libraryBytes, &library)
	if err != nil {
		log.Fatalf("failed to unmarshal library file: %v", err)
	}

	playlists, err := FetchAllPlaylistItems(service)
	if err != nil {
		log.Fatalf("Failed to fetch library: %v", err)
	}
	removedPlaylistItems := []Playlist{}
	notRemovedVideos := []Video{}
	for _, playlist := range playlists {
		removed, notRemoved := partitionPlaylistItems(playlist.PlaylistItems)
		if len(removed) > 0 {
			removedPlaylistItems = append(removedPlaylistItems, Playlist{ID: playlist.ID, Title: playlist.Title, PlaylistItems: removed})
		}
		if len(notRemoved) > 0 {
			for _, playlistItem := range notRemoved {
				notRemovedVideos = append(notRemovedVideos, Video{ID: playlistItem.ID, Title: playlistItem.Title})
			}
		}
	}

	notRemovedVideosJSON, err := json.MarshalIndent(notRemovedVideos, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal library videos: %v", err)
	}
	err = ioutil.WriteFile(LIBRARY_FILE, notRemovedVideosJSON, 0644)
	if err != nil {
		log.Fatalf("failed to write library file: %v", err)
	}

	previouslyRecoveredPlaylistItemsBytes, err := ioutil.ReadFile(RECOVERED_VIDEOS_FILE)
	if err != nil {
		log.Fatalf("failed to read recovered_videos file: %v", err)
	}
	var previouslyRecoveredPlaylistItems []Playlist
	err = json.Unmarshal(previouslyRecoveredPlaylistItemsBytes, &previouslyRecoveredPlaylistItems)
	if err != nil {
		log.Fatalf("failed to unmarshal recovered_videos file: %v", err)
	}
	previouslyRecoveredVideos := []Video{}
	for _, playlist := range previouslyRecoveredPlaylistItems {
		for _, playlistItem := range playlist.PlaylistItems {
			previouslyRecoveredVideos = append(previouslyRecoveredVideos, Video{Title: playlistItem.Title, ID: playlistItem.ID})
		}
	}

	recoveredPlaylistItems := []Playlist{}
	unrecoveredPlaylistItems := []Playlist{}
	fullLibrary := append(library, previouslyRecoveredVideos...)
	for _, playlist := range removedPlaylistItems {
		recoveredItems := []PlaylistItem{}
		unrecoveredItems := []PlaylistItem{}
		for _, playlistItem := range playlist.PlaylistItems {
			recovered := false
			for _, video := range fullLibrary {
				if playlistItem.ID == video.ID {
					playlistItem.Title = video.Title
					recoveredItems = append(recoveredItems, playlistItem)
					recovered = true
					break
				}
			}
			if !recovered {
				unrecoveredItems = append(unrecoveredItems, playlistItem)
			}
		}
		if len(recoveredItems) > 0 {
			recoveredPlaylistItems = append(recoveredPlaylistItems, Playlist{ID: playlist.ID, Title: playlist.Title, PlaylistItems: recoveredItems})
		}
		if len(unrecoveredItems) > 0 {
			unrecoveredPlaylistItems = append(unrecoveredPlaylistItems, Playlist{ID: playlist.ID, Title: playlist.Title, PlaylistItems: unrecoveredItems})
		}
	}

	recoveredPlaylistItemsJSON, err := json.MarshalIndent(recoveredPlaylistItems, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal recovered videos: %v", err)
	}
	err = ioutil.WriteFile(RECOVERED_VIDEOS_FILE, recoveredPlaylistItemsJSON, 0644)
	if err != nil {
		log.Fatalf("failed to write recovered_videos file: %v", err)
	}

	unrecoveredPlaylistItemsJSON, err := json.MarshalIndent(unrecoveredPlaylistItems, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal unrecovered videos: %v", err)
	}
	err = ioutil.WriteFile(UNRECOVERED_VIDEOS_FILE, unrecoveredPlaylistItemsJSON, 0644)
	if err != nil {
		log.Fatalf("failed to write unrecovered_videos file: %v", err)
	}
}
