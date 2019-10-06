package storage

import (
	"blink-backend/database"
	"blink-backend/database/model"
	uuid2 "github.com/google/uuid"
	"os"
)

func craftFilePath(uuid string) string {
	return "./_blink/" + uuid
}

func MakeFile(nickname string, filename string) (string, error) {
	uuid := uuid2.New().String()
	fi, err := os.Create(craftFilePath(uuid))
	if err != nil {
		return "", err
	}
	// close fi on exit and check for its returned error
	defer fi.Close()
	
	db := database.GetInstance().DB
	
	file := &model.File{
		Uuid:     uuid,
		Nickname: nickname,
		Filename: filename,
	}

	db.NewRecord(file)
	db.Create(file)

	return uuid, nil
}

func WriteToFile(uuid string, chunk []byte) error {
	fi, err := os.OpenFile(craftFilePath(uuid), os.O_APPEND | os.O_WRONLY, os.ModeAppend)
	defer fi.Close()
	if err != nil {
		return err
	}
	if _, err := fi.Write(chunk); err != nil {
		return err
	}
	return nil
}
