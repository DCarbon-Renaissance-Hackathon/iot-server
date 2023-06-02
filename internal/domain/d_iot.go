package domain

import "github.com/Dcarbon/iott-cloud/internal/models"

type RIotCreate struct {
	Project  int64             `json:"project" binding:"required"`
	Address  models.EthAddress `json:"address" binding:"required"`
	Type     models.IOTType    `json:"type"  binding:"required"`
	Position *models.Point4326 `json:"position" binding:"required"`
} //@name RIotCreate

type RIotChangeStatus struct {
	IotId  int64                `json:"iotId" form:"iotId" binding:"required"`
	Status *models.DeviceStatus `json:"status" form:"status" binding:"required"`
} //@name RIotChangeStatus

type RIotGetMintSignList struct {
	From  int64 `json:"from" form:"from"  binding:"required"`
	To    int64 `json:"to" form:"to"  binding:"required"`
	IotId int64 `json:"iotId" form:"iotId"  binding:"required"`
} //@name RIotGetMintSignList

type IIot interface {
	Create(*RIotCreate) (*models.IOTDevice, error)
	ChangeStatus(*RIotChangeStatus) (*models.IOTDevice, error)
	GetByBB(min, max *models.Point4326) ([]*models.IOTDevice, error) // boundingbox
	GetIOT(id int64) (*models.IOTDevice, error)
	GetIOTByAddress(addr models.EthAddress) (*models.IOTDevice, error)

	// GetIOTStatus(iotAddr string) models.IOTStatus

	// CreateMetric(*models.Metric) error
	// GetMetrics(iotAddr string, from, to int64) ([]*models.Metric, error)
	// GetRawMetric(metricId string) (*models.Metric, error)

	CreateMint(mint *models.MintSign) error
	GetMintSigns(*RIotGetMintSignList) ([]*models.MintSign, error)
}
