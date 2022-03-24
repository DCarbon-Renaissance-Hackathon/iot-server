package models

type MintSignature struct {
	ID     int64  `gorm:"primary_key"`
	Nonce  int64  `gorm:"index"`
	Amount int64  ``             // Decimal 3
	IoT    string `gorm:"index"` // IoT Address
	R      string ``
	S      string ``
	V      int8   ``
}

func (*MintSignature) TableName() string { return TableNameMintSignature }
