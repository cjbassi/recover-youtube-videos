package utils

import (
	. "github.com/cjbassi/recover-youtube-videos/backend/src/models"
)

func SplitVideos(videos []Video) ([]Video, []Video) {
	removed := make([]Video, 0)
	notRemoved := make([]Video, 0)
	for _, video := range videos {
		if video.Title == "Deleted video" || video.Title == "Private video" {
			removed = append(removed, video)
		} else {
			notRemoved = append(notRemoved, video)
		}
	}
	return removed, notRemoved
}
