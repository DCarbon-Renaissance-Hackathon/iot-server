package models

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
)

// type SensorType int32

// const (
// 	SensorTypeNone           SensorType = 0
// 	SensorTypeFlow           SensorType = 1
// 	SensorTypePower          SensorType = 2
// 	SensorTypeGPS            SensorType = 3
// 	SensorTypeGPSTemperature SensorType = 4
// )

// type DeviceStatus int32

// const (
// 	DeviceStatusReject   DeviceStatus = -1
// 	DeviceStatusRegister DeviceStatus = 0
// 	DeviceStatusSuccess  DeviceStatus = 10
// )

type Sensor struct {
	ID        int64                `json:"id"`
	IotID     int64                `json:"iotId"`
	Address   *EthAddress          `json:"address" gorm:"index:,unique,where:length(address) > 0"`
	Type      dmodels.SensorType   `json:"type"`
	Status    dmodels.DeviceStatus `json:"status"`
	CreatedAt time.Time            `json:"createdAt"`
} // @name Sensor

func (*Sensor) TableName() string { return TableNameSensors }
