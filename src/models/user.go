package models

type User struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
	Token      string `json:"token,omitempty"`
	PictureUrl string `json:"picture_url"`
	CategoryId int    `gorm:"foreignKey:CategoryId" json:"category_id"`
	Role       string `json:"role"`
}
