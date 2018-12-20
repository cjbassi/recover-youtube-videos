package youtube

type PlaylistItem struct {
	ID         string `json:"id"`
	PlaylistID string `json:"playlistId"`
	Title      string `json:"title"`
	Position   int64  `json:"position"`
	Thumbnail  string `json:"thumbnail"`
}

type Playlist struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	PlaylistItems []PlaylistItem `json:"playlistItems"`
}
