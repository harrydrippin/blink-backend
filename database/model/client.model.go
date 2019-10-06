package model

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	Nickname     string `gorm:"PRIMARY_KEY"`
	Latitude     float32
	Longitude    float32
	AwayDistance float32 `gorm:"-"` // ignore
}
