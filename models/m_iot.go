package models

type IOTStatus int

const (
	IOTStatusReject   IOTStatus = -1
	IOTStatusRegister IOTStatus = 0
	IOTStatusSuccess  IOTStatus = 1
)

type IOTType int

const (
	IOTTypeNone         = 0
	IOTTypeDungElectric = 1
)

type IOTDevice struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	Project  int64     `json:"project" `
	Type     int32     `json:"type" `
	Address  string    `json:"address" `
	Status   IOTStatus ``
	Position Point4326 `json:"position" gorm:"column:pos;index;type:geometry(POINT, 4326)"`
}

func (*IOTDevice) TableName() string { return TableNameIOT }
