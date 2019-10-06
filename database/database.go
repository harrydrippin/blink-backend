package database

import (
	"blink-backend/database/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

import (
	"sync"
)

type Database struct {
	DB *gorm.DB
}

var instance *Database
var once sync.Once

func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{}
		instance.Initialize()
	})
	return instance
}

func (database *Database) Initialize() {
	db, err := gorm.Open("sqlite3", "./_blink/database.db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db.AutoMigrate(&model.Client{})
	db.AutoMigrate(&model.File{})
	db.AutoMigrate(&model.Spot{})

	database.DB = db

}

func (database *Database) Close() {
	database.DB.Close()
}
