package models

import (
	"time"
)

type SensorType int32

const (
	SensorTypeNone  SensorType = 0
	SensorTypeFlow  SensorType = 1
	SensorTypePower SensorType = 2
)

type SensorStatus int32

const (
	SensorStatusReject   IOTStatus = -1
	SensorStatusRegister IOTStatus = 0
	SensorStatusSuccess  IOTStatus = 10
)

type Sensor struct {
	ID        int64        ``
	IotID     int64        ``
	Address   EthAddress   ``
	Type      SensorType   `` // CH4, KW, MW, ...
	Status    SensorStatus ``
	CreatedAt time.Time    ``
}

func (*Sensor) TableName() string { return TableNameSensors }

type SensorMetrict struct {
	ID         string    ``
	IotID      int64     ``
	SensorID   int64     ``
	Indicator  float64   ``
	RawMessage string    `` // Hex json of SensorMetricExtract
	CreatedAt  time.Time ``
}

func (*SensorMetrict) TableName() string { return TableNameSensorMetrics }

type SensorMetricExtract struct {
	From      int64   ``
	To        int64   ``
	Indicator float64 ``
	Address   string  `` // Sensor address
}

// type MetrictAggregate struct {
// 	ID        string    ``
// 	IotID     int64     ``
// 	Type      int64     ``
// 	Indicator float64   ``
// 	CreatedAt time.Time ``
// }
