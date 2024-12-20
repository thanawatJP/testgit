package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:100;unique" json:"role_name"`
}
