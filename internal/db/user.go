package db

import (
	"MyLNPU/internal/model"
)

func CreateUser(u *model.User) error {
	return db.Create(u).Error
}

func GetUserByID(openid string) (*model.User, error) {
	var u model.User
	if err := db.First(&u, "open_id = ?", openid).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(u *model.User) error {
	return db.Model(&model.User{}).Where("open_id = ?", u.OpenID).Updates(u).Error
}
