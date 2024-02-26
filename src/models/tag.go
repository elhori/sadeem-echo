package models

type Tag struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pircture_url"`
	IsActive   bool   `json:"is_active"`
}
