package repo

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/Dcarbon/domain"
	"github.com/Dcarbon/libs/dbutils"
	"github.com/Dcarbon/libs/esign"
	"github.com/Dcarbon/libs/utils"
	"github.com/Dcarbon/models"
	"github.com/ethereum/go-ethereum/common"
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

	metric.Extract = models.ExtractMetric{}
	utils.Dump("", metric)

	// err = irepo.CreateMetric(metric)
	// utils.PanicError("Create iot metrics ", err)
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
	var data, err = irepo.GetRawMetric("c419eb47-250e-44ec-98e1-f86b1a813520")
	utils.PanicError("TestGetRawMetrics", err)
	utils.Dump("TestGetRawMetrics", data)
}

func TestCreateMint(t *testing.T) {
	var sign = &models.MintSign{
		Nonce:  1,
		Amount: "0xffaaa1",
		IOT:    iotAddr,
	}
	_, err := sign.Sign(iotPrv)
	utils.PanicError("TestCreateMint", err)

	utils.Dump("MintSign: ", sign)
	// err = irepo.CreateMint(sign)
	// utils.PanicError("TestCreateMint", err)
}

func TestGetMint(t *testing.T) {
	var signeds, err = irepo.GetMintSigns(iotAddr, 0)
	utils.PanicError("TestGetMint", err)
	utils.Dump("TestGetMint", signeds)
}

func TestIsAddress(t *testing.T) {
	isAddr := common.IsHexAddress("0x6cff13d489623029d4d102fa81947527e175ba8d")
	log.Println("Is address: ", isAddr)
}
