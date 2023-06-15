package domain

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type ROpSetStatus struct {
	Id     int64           `json:"-"`      // Iot id
	Status models.OpStatus `json:"status"` // Operator status
}

type RChangeMetric struct {
	IotId    int64             `json:"iotId"`
	SensorId int64             `json:"sensorId"`
	Metric   *models.AllMetric `json:"metric"`
}

type RsGetMetrics struct {
	Id      int64                    `json:"id"`      // Iot id
	Metrics []*models.OpSensorMetric `json:"metrics"` //
}

type IOperator interface {
	SetStatus(req *ROpSetStatus) error
	GetStatus(iotId int64) (*models.OpIotStatus, error)

	ChangeMetrics(*RChangeMetric, dmodels.SensorType) (*models.OpSensorMetric, error)
	GetMetrics(iotId int64) (*RsGetMetrics, error)
}
