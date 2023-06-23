package domain

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

// Identify of sensor. Required id or address
type SensorID struct {
	ID      int64              `json:"id"`
	Address dmodels.EthAddress `json:"address"`
}

type Metric struct {
	ID         string             `json:"id"`
	IotId      int64              `json:"iotId"`
	SensorId   int64              `json:"sensorId"`
	SensorType dmodels.SensorType `json:"sensorType"`
	Indicator  *dmodels.AllMetric `json:"indicator"`
	Data       string             `json:"data"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type RCreateSensor struct {
	IotID   int64              `json:"iotId"`
	Type    dmodels.SensorType `json:"type"`    // CH4, KW, MW, ...
	Address dmodels.EthAddress `json:"address"` // Sensor address
} //@name RCreateSensor

type RChangeSensorStatus struct {
	ID     int64                `json:"id"`     // Sensor id
	Status dmodels.DeviceStatus `json:"status"` //
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
	SensorAddress dmodels.EthAddress `json:"sensorAddress"` //
	Data          string             `json:"data"`          // Hex json of SMExtract
	Signed        string             `json:"signed"`        // Hex of rsv (65bytes)
}

type RCreateSensorMetric struct {
	Data        string             `json:"data" binding:"required"`        // Hex json of SMExtract
	Signed      string             `json:"signed" binding:"required"`      // Hex of rsv (65bytes)
	SignAddress dmodels.EthAddress `json:"signAddress" binding:"required"` //
	IsIotSign   bool               `json:"isIotSign"`                      //
	SensorID    int64              `json:"sensorId" binding:"required"`    //
	IotID       int64              `json:"-"`                              //
} //@name RCreateSensorMetric

type RGetSM struct {
	From     int64 `json:"from" form:"from" binding:"required"`          // Timestamp start
	To       int64 `json:"to" form:"to" binding:"required"`              // Timestamp end
	IotId    int64 `json:"iotId" form:"iotId" binding:"required"`        //
	SensorId int64 `json:"sensorId" form:"sensorId" `                    //
	Skip     int64 `json:"skip" form:"skip"`                             //
	Limit    int64 `json:"limit" form:"limit" binding:"required,max=50"` //
	WithSign bool  `json:"withSign" form:"withSign"`
}

type RSMAggregate struct {
	From     int64 `json:"from" form:"from" binding:"required"`         // Timestamp start
	To       int64 `json:"to" form:"to" binding:"required"`             // Timestamp end
	IotId    int64 `json:"iotId" form:"iotId" binding:"required"`       //
	SensorId int64 `json:"sensorId" form:"sensorId" binding:"required"` //
	Interval int   `json:"interval" form:"interval" binding:"required"` // 1 : day 2: month
}

type TimeValue struct {
	Time time.Time `json:"time"`
	Val  float64   `json:"value"`
} // @name TimeValue

type ISensor interface {
	SetOperatorCache(op IOperator)

	CreateSensor(*RCreateSensor) (*models.Sensor, error)
	ChangeSensorStatus(*RChangeSensorStatus) (*models.Sensor, error)
	GetSensor(*SensorID) (*models.Sensor, error)
	GetSensors(*RGetSensors) ([]*models.Sensor, error)
	GetSensorType(req *SensorID) (dmodels.SensorType, error)

	CreateSM(*RCreateSM) (*models.SmSignature, error)                     // old
	CreateSensorMetric(*RCreateSensorMetric) (*models.SmSignature, error) // New

	GetMetrics(*RGetSM) ([]*Metric, error)
	GetAggregatedMetrics(*RSMAggregate) ([]*TimeValue, error)
}
