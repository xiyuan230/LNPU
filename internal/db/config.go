package db

import (
	"MyLNPU/internal/model"
)

func GetSystemNotice() (*model.Notice, error) {
	var c model.Notice
	if err := db.First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateSystemNotice(notice *model.Notice) error {
	return db.Save(notice).Error
}
