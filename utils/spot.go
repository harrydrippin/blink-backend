package utils

import (
	"blink-backend/database"
	"blink-backend/database/model"
	"github.com/jinzhu/gorm"
)

func MakeSpot(nickname string, lat float32, lng float32, uuid string) uint32 {
	db := database.GetInstance().DB
	spot := model.Spot{
		Uuid:      uuid,
		Nickname:  nickname,
		Latitude:  lat,
		Longitude: lng,
		HitCount:  0,
	}
	db.Create(&spot)

	return uint32(spot.ID)
}

func GetSpotById(id uint32) model.Spot {
	db := database.GetInstance().DB
	var spot model.Spot
	db.Where("id = ?", id).First(&spot)
	return spot
}

func GetSpotsByNickname(nickname string) []model.Spot {
	db := database.GetInstance().DB
	var spots []model.Spot
	db.Where("nickname = ?", nickname).Find(&spots)
	return spots
}

func GetSpotsByLocation(lat float32, lng float32) []model.Spot {
	db := database.GetInstance().DB
	var spots []model.Spot
	db.Where(
		"latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ?",
		lat+0.004, lat-0.004,
		lng+0.004, lng-0.004,
	).Find(&spots)

	return spots
}

func GetFileFromSpot(id uint32) string {
	db := database.GetInstance().DB
	var spot model.Spot
	if err := db.Where("id = ?", id).First(&spot).Error; gorm.IsRecordNotFoundError(err) {
		return ""
	}

	var file model.File
	if err := db.Where("uuid = ?", spot.Uuid).First(&file).Error; gorm.IsRecordNotFoundError(err) {
		return ""
	}

	return file.CraftFileLink()
}
