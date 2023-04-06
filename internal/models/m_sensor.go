package models

import (
	"time"
)

type SensorType int32

const (
	SensorTypeNone  SensorType = 0
	SensorTypeFlow  SensorType = 1
	SensorTypePower SensorType = 2
	SensorTypeGPS   SensorType = 3
)

type SensorStatus int32

const (
	SensorStatusReject   SensorStatus = -1
	SensorStatusRegister SensorStatus = 0
	SensorStatusSuccess  SensorStatus = 10
)

type Sensor struct {
	ID        int64        ``
	IotID     int64        ``
	Address   *EthAddress  `gorm:"index:,unique,where:length(address) > 0"`
	Type      SensorType   ``
	Status    SensorStatus ``
	CreatedAt time.Time    ``
}

func (*Sensor) TableName() string { return TableNameSensors }
