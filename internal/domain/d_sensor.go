package domain

import (
	"time"

	"github.com/Dcarbon/iott-cloud/internal/models"
)

type RCreateSensor struct {
	Type      string            `json:"type"`      // CH4, KW, MW, ...
	Address   models.EthAddress `json:"address"`   // Sensor address
	CreatedAt time.Time         `json:"createdAt"` //
}

type RChangeSensorStatus struct {
	ID int64
}

type RGetSensor struct {
}

type RGetSensors struct {
}

type RCreateSensorMetric struct {
}

type RGetMetrics struct {
}

type ISensor interface {
	CreateSensor(*RCreateSensor) (*models.Sensor, error)
	ChangeStatus(*RChangeSensorStatus) (*models.Sensor, error)
	GetSensor(*RGetSensor) (*models.Sensor, error)
	GetSensors(*RGetSensors) ([]*models.Sensor, error)

	CreateMetric(*RCreateSensorMetric) (*models.SensorMetrict, error)
	GetMetric(*RGetMetrics) ([]*models.SensorMetrict, error)
}
