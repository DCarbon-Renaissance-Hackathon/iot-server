package domain

import "github.com/Dcarbon/iott-cloud/internal/models"

type IIot interface {
	Create(iot *models.IOTDevice) error
	ChangeStatus(iotAddr string, status models.IOTStatus) (*models.IOTDevice, error)
	GetByBB(min, max *models.Point4326) ([]*models.IOTDevice, error) // boundingbox
	GetIOT(id int64) (*models.IOTDevice, error)
	GetIOTByAddress(addr models.EthAddress) (*models.IOTDevice, error)

	// GetIOTStatus(iotAddr string) models.IOTStatus

	// CreateMetric(*models.Metric) error
	// GetMetrics(iotAddr string, from, to int64) ([]*models.Metric, error)
	// GetRawMetric(metricId string) (*models.Metric, error)

	CreateMint(mint *models.MintSign) error
	GetMintSigns(iotAddr string, fromNonce int) ([]*models.MintSign, error)
}
