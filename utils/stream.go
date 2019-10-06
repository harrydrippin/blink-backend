package utils

import (
	"blink-backend/blink"
	"blink-backend/database"
	"blink-backend/database/model"
)

func SetReceiverInfo(location *blink.Location, nickname string) error {
	db := database.GetInstance().DB

	var client model.Client
	db.Where("nickname = ?", nickname).First(&client)
	client.Latitude = location.GetLatitude()
	client.Longitude = location.GetLongitude()
	db.Save(&client)

	return nil
}
