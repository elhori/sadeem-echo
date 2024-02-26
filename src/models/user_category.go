package models

type UserCategory struct {
	Id   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
