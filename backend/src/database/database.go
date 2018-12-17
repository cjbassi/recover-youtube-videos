package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"

	. "github.com/cjbassi/recover-youtube-videos/backend/src/models"
)

type Database struct {
	logger     *logrus.Entry
	Connection *gorm.DB
}

func Setup(logger *logrus.Entry, databaseURL string) (*Database, error) {
	db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	database := &Database{
		logger:     logger,
		Connection: db,
	}
	return database, nil
}

func (db *Database) Close() {
	db.Connection.Close()
}

func (db *Database) HardMigrate() {
	db.logger.Infof("Dropping tables if they exist...")

	db.Connection.DropTableIfExists(&Video{}, &Playlist{}, &Channel{})

	db.logger.Infof("Creating tables...")

	db.Connection.CreateTable(&Channel{}, &Playlist{}, &Video{})

	db.Connection.Model(&Video{}).AddForeignKey("playlist_id", "playlists(id)", "RESTRICT", "RESTRICT")
	db.Connection.Model(&Playlist{}).AddForeignKey("channel_id", "channels(id)", "RESTRICT", "RESTRICT")

	db.logger.Infof("Created tables")
}

func (db *Database) SoftMigrate() {
	db.logger.Infof("Automigrating tables")

	db.Connection.AutoMigrate(&Channel{}, &Playlist{}, &Video{})
}
