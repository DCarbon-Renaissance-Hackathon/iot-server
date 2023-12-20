package models

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
)

const (
	TableNameXM = "x_metric"
)

// Experiment sensor metric
type XSMetric struct {
	Id         string             `json:"id,omitempty"           `
	IotAddress dmodels.EthAddress `json:"iotAddress,omitempty"   gorm:"index_ca,priority:2"`
	SensorType dmodels.SensorType `json:"sensorType,omitempty"   `
	Metric     *dmodels.AllMetric `json:"metric,omitempty"       gorm:"type:json"`
	CreatedAt  time.Time          `json:"createdAt,omitempty"    gorm:"index_ca,priority:1"`
}

func (*XSMetric) TableName() string { return TableNameXM }
