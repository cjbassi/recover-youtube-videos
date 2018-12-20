package database

type Video struct {
	ID        string `json:"id" gorm:"primary_key"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}
