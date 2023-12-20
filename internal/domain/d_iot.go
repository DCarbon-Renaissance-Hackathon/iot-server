package domain

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type Sort int

const (
	SortASC  Sort = 0
	SortDesc Sort = 1
)

type RIotCreate struct {
	Project  int64              `json:"project" binding:"required"`
	Address  dmodels.EthAddress `json:"address" binding:"required"`
	Type     models.IOTType     `json:"type"  binding:"required"`
	Position *models.Point4326  `json:"position" binding:"required"`
}

type RIotChangeStatus struct {
	IotId  int64                 `json:"iotId" form:"iotId" binding:"required"`
	Status *dmodels.DeviceStatus `json:"status" form:"status" binding:"required"`
} //@name RIotChangeStatus

type RIotUpdate struct {
	IotId    int64             `json:"iotId" form:"iotId" binding:"required"`
	Position *models.Point4326 `json:"position" binding:"required"`
} //@name RIotChangeStatus

type RIotGetList struct {
	Skip      int                  `json:"skip" form:"skip"`
	Limit     int                  `json:"limit" form:"limit" binding:"max=50"`
	ProjectId int64                `json:"projectId" form:"projectId" binding:"required"`
	Status    dmodels.DeviceStatus `json:"status" form:"status"`
}

type RIotMint struct {
	Nonce  int64  `json:"nonce" binding:"required"`  //
	Amount string `json:"amount" binding:"required"` // Hex
	Iot    string `json:"iot" binding:"required"`    // IoT Address
	R      string `json:"r" binding:"required"`      //
	S      string `json:"s" binding:"required"`      //
	V      string `json:"v" binding:"required"`      //
} // @name RIotMint

type RIotGetMintSignList struct {
	From  int64 `json:"from" form:"from" binding:"required"`
	To    int64 `json:"to" form:"to" binding:""`
	IotId int64 `json:"iotId" uri:"iotId" binding:"required"`
	Sort  Sort  `json:"sort" form:"sort"`
	Limit int   `json:"limit" form:"limit"`
} //@name RIotGetMintSignList

type RIotGetMintedList struct {
	From     int64 `json:"from" form:"from" binding:"required"`
	To       int64 `json:"to" form:"to" binding:"required"`
	IotId    int64 `json:"iotId" form:"iotId" binding:""`
	Interval int   `json:"interval" form:"interval"` // 1 : day 2: month
} //@name RIotGetMintedList

type RIotCount struct {
} //@name RIotCount

type RIsIotActiced struct {
	From  int64 `json:"from" form:"from" binding:"required"`
	To    int64 `json:"to" form:"to" binding:"required"`
	IotId int64 `json:"iotId" form:"iotId" binding:"required"`
} //@name RIsIotActiced

type PositionId struct {
	Id       int64          `json:"id"`
	Position *dmodels.Coord `json:"position"`
} //@name PositionId

type IIot interface {
	Create(*RIotCreate) (*models.IOTDevice, error)
	Update(req *RIotUpdate) (*models.IOTDevice, error)
	ChangeStatus(*RIotChangeStatus) (*models.IOTDevice, error)
	GetIot(id int64) (*models.IOTDevice, error)
	GetIots(*RIotGetList) ([]*models.IOTDevice, error)
	GetIotPositions(*RIotGetList) ([]*PositionId, error)

	GetIotByAddress(addr dmodels.EthAddress) (*models.IOTDevice, error)

	// GetIOTStatus(iotAddr string) models.IOTStatus

	// CreateMetric(*models.Metric) error
	// GetMetrics(iotAddr string, from, to int64) ([]*models.Metric, error)
	// GetRawMetric(metricId string) (*models.Metric, error)

	CreateMint(mint *RIotMint) error
	GetMintSigns(*RIotGetMintSignList) ([]*models.MintSign, error)
	GetMinted(*RIotGetMintedList) ([]*models.Minted, error)

	CountIot(*RIotCount) (int64, error)
	IsIotActived(req *RIsIotActiced) (bool, error)
}
