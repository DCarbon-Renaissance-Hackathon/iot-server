package domain

import "github.com/Dcarbon/models"

type IIot interface {
	Create(iot *models.IOTDevice) error
	ChangeStatus(iotAddr string, status models.IOTStatus) (*models.IOTDevice, error)
	GetByBB(min, max *models.Point4326) ([]*models.IOTDevice, error) // boundingbox
	// GetIOTStatus(iotAddr string) models.IOTStatus

	CreateMetric(*models.Metric) error
	GetMetrics(iotAddr string, from, to int64) ([]*models.Metric, error)
	GetRawMetric(iotAddr string, metricId string) (*models.Metric, error)
}
