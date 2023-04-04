package models

const (
	RoleAdmin = "admin"
)

type User struct {
	ID      int64      `json:"id" gorm:"primary_key"`
	Role    string     `json:"role"`
	Name    string     `json:"name" `
	Address EthAddress `json:"address" gorm:"unique"` // Eth address
}

func (*User) TableName() string { return TableNameUser }
