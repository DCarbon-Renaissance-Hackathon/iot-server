package domain

import (
	"github.com/Dcarbon/iott-cloud/internal/models"
)

// Identify of sensor. Required id or address
type SensorID struct {
	ID      int64             `json:"id"`
	Address models.EthAddress `json:"address"`
}

type RCreateSensor struct {
	IotID   int64             `json:"iotId"`
	Type    models.SensorType `json:"type"`    // CH4, KW, MW, ...
	Address models.EthAddress `json:"address"` // Sensor address
}

type RChangeSensorStatus struct {
	ID     int64               `json:"id"`
	Status models.SensorStatus `json:"status"`
}

// type RGetSensor struct {
// 	ID      int64             `json:"id"`
// 	Address models.EthAddress `json:"address"`
// }

type RGetSensors struct {
	Skip  int   `json:"skip"`  //
	Limit int   `json:"limit"` //
	IotId int64 `json:"iotId"` //
}

type RCreateSM struct {
	SensorAddress models.EthAddress `json:"sensorAddress"` //
	Data          string            `json:"data"`          // Hex json of SMExtract
	Signed        string            `json:"signed"`        // Hex of rsv (65bytes)
}

type RCreateSMFromIOT struct {
	Data       string            `json:"data"`     // Hex json of SMExtract
	Signed     string            `json:"signed"`   // Hex of rsv (65bytes)
	SensorID   int64             `json:"sensorId"` //
	IotAddress models.EthAddress `json:"iot"`      //
	IotID      int64             `json:"iotId"`    //
}

type RGetSM struct {
	From  int64 `json:"from"`
	To    int64 `json:"to"`
	IotId int64 `json:"iotId"`
}

type ISensor interface {
	CreateSensor(*RCreateSensor) (*models.Sensor, error)
	ChangeSensorStatus(*RChangeSensorStatus) (*models.Sensor, error)
	GetSensor(*SensorID) (*models.Sensor, error)
	GetSensors(*RGetSensors) ([]*models.Sensor, error)

	CreateSM(*RCreateSM) (*models.SmSignature, error)
	CreateSMFromIot(*RCreateSMFromIOT) (*models.SmSignature, error)

	GetMetrics(*RGetSM) ([]*models.SmFloat, error)
}
