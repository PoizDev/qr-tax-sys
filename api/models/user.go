package models

type User struct {
	UserID   uint   `json:"id" gorm:"primaryKey;autoIncrement;column:user_id"`
	Username string `json:"username" gorm:"not null;"`
	Password string `json:"password" gorm:"not null;"`
	Email    string `json:"email" gorm:"not null;"`
}
