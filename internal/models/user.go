package models

import "time"

const (
	RoleAdmin = "admin"
)

type UserType int

const (
	UserTypePersonal UserType = 1
	UserTypeCompany  UserType = 2
	UserTypeFund     UserType = 3
)

type User struct {
	ID        int64      `json:"id" gorm:"primary_key"` //
	Role      string     `json:"role"`                  //
	Name      string     `json:"name" `                 //
	Address   EthAddress `json:"address" gorm:"unique"` // Eth address
	TaxCode   string     `json:"taxCode"`
	Phone     string     `json:"phone"`
	Type      UserType   `json:"type"`
	CreatedAt time.Time  `json:"-"`
}

func (*User) TableName() string { return TableNameUser }
