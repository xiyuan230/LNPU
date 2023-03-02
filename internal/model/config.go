package model

type Notice struct {
	ID     int    `json:"id" gorm:"primaryKey" binding:"required"`
	Notice string `json:"notice" binding:"required"`
}
