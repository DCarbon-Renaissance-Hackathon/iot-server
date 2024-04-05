package domain

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

type RXSMCreate struct {
	Address    dmodels.EthAddress `json:"address"`
	SensorType dmodels.SensorType `json:"sensorType"`
	Metric     *dmodels.AllMetric `json:"metric"`
} // @name RXSMCreate

type RXSMGetList struct {
	Skip       int                `json:"skip"        form:"skip"`
	Limit      int                `json:"limit"       form:"limit"`
	From       int64              `json:"from"        form:"from"`
	To         int64              `json:"to"          form:"to"`
	Address    dmodels.EthAddress `json:"address"     form:"address"`
	SensorType dmodels.SensorType `json:"sensorType"  form:"sensorType"`
} // @name RXSMGetList

type IXSM interface {
	Create(*RXSMCreate) (*models.XSMetric, error)
	GetList(*RXSMGetList) ([]*models.XSMetric, error)
}
