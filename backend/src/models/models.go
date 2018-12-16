package models

type Video struct {
	ID         string `json:"id" gorm:"primary_key"`
	Title      string `json:"title"`
	PlaylistID string `json:playlist_id gorm:"primary_key"`
}

type Playlist struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Channel   *Channel `json:"channel" gorm:"foreignkey:ChannelID"`
	ChannelID string   `json:"channel_id"`
	Videos    []Video  `json:"videos" gorm:"foreignkey:PlaylistID"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ChannelID string `json:"channel_id"`
}
