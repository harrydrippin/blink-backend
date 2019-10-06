package utils

import (
	"blink-backend/database"
	"blink-backend/database/model"
	"github.com/jinzhu/gorm"
)

func CheckNickname(nickname string) (bool, error) {
	db := database.GetInstance().DB
	var client model.Client
	if err := db.Where("Nickname = ?", nickname).First(&client).Error; gorm.IsRecordNotFoundError(err) {
		return true, nil
	}
	return false, nil
}

func SubmitNickname(nickname string) (bool, error) {
	checkResult, err := CheckNickname(nickname)
	if err != nil || !checkResult {
		return false, err
	}

	client := model.Client{
		Nickname:     nickname,
		Latitude:     0.0,
		Longitude:    0.0,
		AwayDistance: 0.0,
	}

	db := database.GetInstance().DB
	db.Create(&client)

	return true, nil
}
