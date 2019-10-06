package model

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Uuid     string `gorm:"PRIMARY_KEY"`
	Nickname string
	Filename string
}

func (f *File) CraftFileLink() string {
	return "/" + f.Uuid + "?filename=" + f.Filename
}
