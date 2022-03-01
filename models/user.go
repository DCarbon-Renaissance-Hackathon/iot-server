package models

type User struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	Role     string `json:"role"`
	Name     string `json:"name" `
	EAddress string `json:"eaddress" `
}

func (*User) TableName() string { return TableNameUser }
