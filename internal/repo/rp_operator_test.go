package repo

import (
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

var opTest *OperatorRepo

func init() {
	var err error
	opTest, err = NewOperatorRepo()
	utils.PanicError("", err)
}

func TestSetStatus(t *testing.T) {
	err := opTest.SetStatus(&domain.ROpSetStatus{
		Id:     1,
		Status: models.OpStatusWarning,
	})
	utils.PanicError("", err)
}

func TestGetStatus(t *testing.T) {
	status, err := opTest.GetStatus(1)
	utils.PanicError("TestGetStatus", err)
	utils.Dump("Status", status)
}

func TestChangeMetricsGPS(t *testing.T) {
	data, err := opTest.ChangeMetrics(
		&domain.RChangeMetric{
			IotId:    1,
			SensorId: 1,
			Metric: &models.AllMetric{
				GPSMetric: models.GPSMetric{
					Lat: 1,
					Lng: 1,
				},
			},
		},
		models.SensorTypeGPS)
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestChangeMetricsPower(t *testing.T) {
	data, err := opTest.ChangeMetrics(
		&domain.RChangeMetric{
			IotId:    1,
			SensorId: 2,
			Metric: &models.AllMetric{
				DefaultMetric: models.DefaultMetric{
					Val: 0.5,
				},
			},
		},
		models.SensorTypePower,
	)
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestGetMetrics(t *testing.T) {
	metrics, err := opTest.GetMetrics(1)
	utils.PanicError("", err)
	utils.Dump("Metrics: ", metrics)
}
