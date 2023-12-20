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
	Skip       int                `json:"skip"`
	Limit      int                `json:"limit"`
	From       int64              `json:"from"`
	To         int64              `json:"to"`
	Address    dmodels.EthAddress `json:"address"`
	SensorType dmodels.SensorType `json:"sensorType"`
} // @name RXSMGetList

type IXSM interface {
	Create(*RXSMCreate) (*models.XSMetric, error)
	GetList(*RXSMGetList) ([]*models.XSMetric, error)
}
