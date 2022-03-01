package repo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Dcarbon/domain"
	"github.com/Dcarbon/libs/dbutils"
	"github.com/Dcarbon/libs/esign"
	"github.com/Dcarbon/libs/utils"
	"github.com/Dcarbon/models"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var irepo = mustIOTRepo()
var iotPrv = utils.StringEnv("IOT_PRIVATE", "")
var iotAddr = utils.StringEnv("IOT_ADDRESS", "")

func mustIOTRepo() domain.IIot {
	irp, err := NewIOTRepo(dbUrl)
	if nil != err {
		panic(err.Error())
	}
	return irp
}

func TestIOTCreate(t *testing.T) {
	err := irepo.Create(&models.IOTDevice{
		Project: 0,
		Type:    models.IOTTypeDungElectric,
		Address: iotAddr,
		Position: models.Point4326{
			Lat: 21.015462,
			Lng: 105.804904,
		},
	})

	utils.PanicError("Create iot device ", err)
}

func TestIOTChangeStatus(t *testing.T) {
	var data, err = irepo.ChangeStatus(
		"0x1064F6f232bdD6E38a248C0C3a1456b023f05e3B",
		models.IOTStatusSuccess,
	)
	utils.PanicError("Update iot status ", err)
	utils.Dump("TestIOTChangeStatus", data)
}

func TestIOTGetByBB(t *testing.T) {
	var data, err = irepo.GetByBB(
		&models.Point4326{Lng: 104.1, Lat: 20},
		&models.Point4326{Lng: 106.1, Lat: 22},
	)
	utils.PanicError("TestIOTGetByBB", err)
	utils.Dump("TestIOTGetByBB", data)
}

func TestIOTCreateMetrics(t *testing.T) {
	var from = time.Now().Unix() - 100
	var to = from + 100

	var extract = &models.ExtractMetric{
		From: from,
		To:   to,
		Position: models.Point4326{
			Lat: 21.015462,
			Lng: 105.804904,
		},
		Metrics: dbutils.MapSFloat{
			"CH4": 100,
			"N2O": 100,
		},
	}
	rawExtract, err := json.Marshal(extract)
	utils.PanicError("Marshal extract ", err)

	rawSigned, err := esign.SignPersonal(iotPrv, rawExtract)
	utils.PanicError("Sign extract ", err)

	var metric = &models.Metric{
		Address:   "0x6CFF13d489623029d4d102Fa81947527E175BA8D",
		Data:      hexutil.Encode(rawExtract),
		Signed:    hexutil.Encode(rawSigned),
		Extract:   *extract,
		CreatedAt: time.Now(),
	}

	err = irepo.CreateMetric(metric)
	utils.PanicError("Create iot metrics ", err)
}

func TestGetMetrics(t *testing.T) {
	var now = time.Now().Unix()
	var from = now - 86400*365*2
	var data, err = irepo.GetMetrics(
		"0x6CFF13d489623029d4d102Fa81947527E175BA8D",
		from,
		now,
	)
	utils.PanicError("TestGetMetrics", err)
	utils.Dump("TestGetMetrics", data)
}

func TestGetRawMetrics(t *testing.T) {
	var data, err = irepo.GetRawMetric(
		"0x6CFF13d489623029d4d102Fa81947527E175BA8D",
		"c419eb47-250e-44ec-98e1-f86b1a813520",
	)
	utils.PanicError("TestGetRawMetrics", err)
	utils.Dump("TestGetRawMetrics", data)
}
