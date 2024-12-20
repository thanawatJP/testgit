package models

type UserAuth struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"size:100" json:"first_name"`
	LastName  string `gorm:"size:100" json:"last_name"`
	Email     string `gorm:"size:255;unique" json:"email"`
	Password  string `gorm:"size:255" json:"password"`
	RoleID    uint   `gorm:"column:role_id" json:"role_id"`
	Role      Role   `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}
