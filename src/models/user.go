package models

type User struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	PictureUrl string `json:"pircture_url"`
	CategoryId int    `gorm:"foreignKey:CategoryId" json:"category_id"`
}
