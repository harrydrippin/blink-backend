package model

import (
	"github.com/jinzhu/gorm"
)

type Spot struct {
	gorm.Model
	Uuid      string
	Nickname  string
	Latitude  float32
	Longitude float32
	HitCount  int
}
