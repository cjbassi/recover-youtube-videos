package utils

import (
	"github.com/cjbassi/recover-youtube-videos/backend/src/api/youtube"
)

func SplitPlaylistItems(items []youtube.PlaylistItem) ([]youtube.PlaylistItem, []youtube.PlaylistItem) {
	removed := make([]youtube.PlaylistItem, 0)
	notRemoved := make([]youtube.PlaylistItem, 0)
	for _, item := range items {
		if item.Title == "Deleted video" || item.Title == "Private video" {
			removed = append(removed, item)
		} else {
			notRemoved = append(notRemoved, item)
		}
	}
	return removed, notRemoved
}
