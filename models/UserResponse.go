package models

type UserResponse struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	Response bool `json:"response"`
	Age      int  `json:"age"`
}
