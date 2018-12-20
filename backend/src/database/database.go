package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	*gorm.DB
}

func Setup(dbURI string) (*DB, error) {
	_db, err := gorm.Open("postgres", dbURI)
	return &DB{_db}, err
}

func (db *DB) HardMigrate() {
	log.Printf("Dropping tables if they exist...")

	db.DropTableIfExists(&Video{})

	log.Printf("Creating tables...")

	db.CreateTable(&Video{})

	log.Printf("Created tables")
}

func (db *DB) SoftMigrate() {
	log.Printf("Automigrating tables")
	db.AutoMigrate(&Video{})
}
