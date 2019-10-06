package utils

import (
	"blink-backend/database"
	"blink-backend/database/model"
)

func GetClientsByLocation(lat float32, lng float32) []model.Client {
	db := database.GetInstance().DB
	var clients []model.Client
	db.Where(
		"latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ?",
		lat+0.004, lat-0.004,
		lng+0.004, lng-0.004,
	).Find(&clients)

	return clients
}

func GetClientsByName(nickname string) []model.Client {
	db := database.GetInstance().DB
	var clients []model.Client
	db.Where(
		"nickname = ?", nickname,
	).Find(&clients)

	return clients
}
