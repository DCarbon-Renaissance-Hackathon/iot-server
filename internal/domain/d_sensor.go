package domain

import (
	"time"

	"github.com/Dcarbon/iott-cloud/internal/models"
)

// Identify of sensor. Required id or address
type SensorID struct {
	ID      int64             `json:"id"`
	Address models.EthAddress `json:"address"`
}

type Metric struct {
	ID        string            `json:"id"`
	IotId     int64             `json:"iotId"`
	SensorId  int64             `json:"sensorId"`
	Indicator *models.AllMetric `json:"indicator"`
	CreatedAt time.Time         `json:"createdAt"`
}

type RCreateSensor struct {
	IotID   int64             `json:"iotId"`
	Type    models.SensorType `json:"type"`    // CH4, KW, MW, ...
	Address models.EthAddress `json:"address"` // Sensor address
} //@name RCreateSensor

type RChangeSensorStatus struct {
	ID     int64               `json:"id"`     // Sensor id
	Status models.DeviceStatus `json:"status"` //
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
	From  int64 `json:"from"` // Timestamp start
	To    int64 `json:"to"`   // Timestamp end
	IotId int64 `json:"iotId"`
}

type ISensor interface {
	SetOperatorCache(op IOperator)

	CreateSensor(*RCreateSensor) (*models.Sensor, error)
	ChangeSensorStatus(*RChangeSensorStatus) (*models.Sensor, error)
	GetSensor(*SensorID) (*models.Sensor, error)
	GetSensors(*RGetSensors) ([]*models.Sensor, error)
	GetSensorType(req *SensorID) (models.SensorType, error)

	CreateSM(*RCreateSM) (*models.SmSignature, error)
	CreateSMFromIot(*RCreateSMFromIOT) (*models.SmSignature, error)

	GetMetrics(*RGetSM) ([]*Metric, error)
}
