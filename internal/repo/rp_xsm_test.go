package repo

import (
	"testing"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/rss"
)

var xsmTest *XSMImpl

func TestCreate(t *testing.T) {
	var req = &domain.RXSMCreate{
		Address:    "0xE445517AbB524002Bb04C96F96aBb87b8B19b53d",
		SensorType: 1,
		Metric: &dmodels.AllMetric{
			DefaultMetric: dmodels.DefaultMetric{
				Val: 122.0,
			},
		},
	}
	var data, err = getXSMTest().Create(req)
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestGetList(t *testing.T) {
	var req = &domain.RXSMGetList{
		Skip:       0,
		Limit:      10,
		Address:    "0xe445517abb524002bb04c96f96abb87b8b19b53d",
		SensorType: 1,
	}
	var data, err = getXSMTest().GetList(req)
	utils.PanicError("", err)
	utils.Dump("", data)
}

func getXSMTest() *XSMImpl {
	if xsmTest != nil {
		return xsmTest
	}

	var db = rss.GetDB()
	var err error
	xsmTest, err = NewXSMImpl(db)
	utils.PanicError("Create XSM test", err)
	return xsmTest
}
