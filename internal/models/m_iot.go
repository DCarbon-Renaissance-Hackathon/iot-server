package models

type IOTStatus int

const (
	IOTStatusReject   IOTStatus = -1
	IOTStatusRegister IOTStatus = 0
	IOTStatusSuccess  IOTStatus = 10
)

type IOTType int

const (
	IOTTypeNone        IOTType = 0
	IOTTypeWindPower   IOTType = 10
	IOTTypeSolarPower  IOTType = 11
	IOTTypeBurnMethane IOTType = 20
	IOTTypeFertilizer  IOTType = 30
	IOTTypeTrash       IOTType = 31
)

type IOTDevice struct {
	ID       int64      `json:"id" gorm:"primary_key"`
	Project  int64      `json:"project" `
	Address  EthAddress `json:"address" gorm:"unique"`
	Type     IOTType    `json:"type" `
	Status   IOTStatus  `json:"status"`
	Position Point4326  `json:"position" gorm:"type:geometry(POINT, 4326)"`
}

func (*IOTDevice) TableName() string { return TableNameIOT }
