package repo

import (
	"log"
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

var iotRepoTest domain.IIot
var iotPrv = utils.StringEnv("IOT_PRIVATE", "")
var iotAddr = utils.StringEnv("IOT_ADDRESS", "")

var testDomainMinter = esign.MustNewERC712(
	&esign.TypedDataDomain{
		Name:              "CARBON",
		Version:           "1",
		ChainId:           1337,
		VerifyingContract: "0x7BDDCb9699a3823b8B27158BEBaBDE6431152a85",
	},
	esign.MustNewTypedDataField(
		"Mint",
		esign.TypedDataStruct,
		esign.MustNewTypedDataField("iot", esign.TypedDataAddress),
		esign.MustNewTypedDataField("amount", "uint256"),
		esign.MustNewTypedDataField("nonce", "uint256"),
	),
)

func init() {
	var err error
	iotRepoTest, err = NewIOTRepo(testDomainMinter)
	if nil != err {
		panic(err.Error())
	}
}

func TestIOTCreate(t *testing.T) {
	// 21.016975, 105.780917
	iot, err := iotRepoTest.Create(&domain.RIotCreate{
		Project: 1,
		Type:    models.IOTTypeBurnMethane,
		Address: "0xe445517abb524002bb04c96f96abb87b8b19b53d",
		Position: &models.Point4326{
			Lat: 21.016975,
			Lng: 105.780917,
		},
	})

	utils.PanicError("Create iot device ", err)
	utils.Dump("IOT created", iot)
}

func TestIOTChangeStatus(t *testing.T) {
	var status = dmodels.DeviceStatusSuccess
	var data, err = iotRepoTest.ChangeStatus(
		&domain.RIotChangeStatus{
			IotId:  292,
			Status: &status,
		},
	)
	utils.PanicError("Update iot status ", err)
	utils.Dump("TestIOTChangeStatus", data)
}

func TestIotUpdate(t *testing.T) {
	var data, err = iotRepoTest.Update(
		&domain.RIotUpdate{
			IotId: 294,
			Position: &models.Point4326{
				Lng: 105.553486,
				Lat: 21.706776,
			},
		},
	)
	utils.PanicError("Update iot status ", err)
	utils.Dump("TestIOTChangeStatus", data)
}

func TestIOTGetIOT(t *testing.T) {
	var data, err = iotRepoTest.GetIot(1)
	utils.PanicError("TestIOTGetIOT", err)
	utils.Dump("TestIOTGetIOT", data)
}

func TestIOTGetIOTPosition(t *testing.T) {
	var data, err = iotRepoTest.GetIotPositions(&domain.RIotGetList{
		Status: dmodels.DeviceStatusSuccess,
	})
	utils.PanicError("TestIOTGetIOTPosition", err)
	utils.Dump("TestIOTGetIOTPosition", data)
}

func TestIOTGetIOTByAddress(t *testing.T) {
	var data, err = iotRepoTest.GetIotByAddress(
		dmodels.EthAddress("0x72ef9da2af1d657b3fd16e93fb9e6d82c4c615f1"),
	)
	utils.PanicError("TestIOTGetIOTPosition", err)
	utils.Dump("TestIOTGetIOTPosition", data)
}

// func TestIOTCreateMetrics(t *testing.T) {
// 	var from = time.Now().Unix() - 100
// 	var to = from + 100

// 	var extract = &models.ExtractMetric{
// 		From: from,
// 		To:   to,
// 		Position: models.Point4326{
// 			Lat: 21.015462,
// 			Lng: 105.804904,
// 		},
// 		Metrics: dbutils.MapSFloat{
// 			"CH4": 100,
// 			"N2O": 100,
// 		},
// 	}
// 	rawExtract, err := json.Marshal(extract)
// 	utils.PanicError("Marshal extract ", err)

// 	rawSigned, err := esign.SignPersonal(iotPrv, rawExtract)
// 	utils.PanicError("Sign extract ", err)

// 	var metric = &models.Metric{
// 		Address:   "0x6CFF13d489623029d4d102Fa81947527E175BA8D",
// 		Data:      hexutil.Encode(rawExtract),
// 		Signed:    hexutil.Encode(rawSigned),
// 		Extract:   *extract,
// 		CreatedAt: time.Now(),
// 	}

// 	metric.Extract = models.ExtractMetric{}
// 	utils.Dump("", metric)

// 	err = iotRepoTest.CreateMetric(metric)
// 	utils.PanicError("Create iot metrics ", err)
// }

// func TestGetMetrics(t *testing.T) {
// 	var now = time.Now().Unix()
// 	var from = now - 86400*365*2
// 	var data, err = iotRepoTest.GetMetrics(
// 		"0x6CFF13d489623029d4d102Fa81947527E175BA8D",
// 		from,
// 		now,
// 	)
// 	utils.PanicError("TestGetMetrics", err)
// 	utils.Dump("TestGetMetrics", data)
// }

// func TestGetRawMetrics(t *testing.T) {
// 	var data, err = iotRepoTest.GetRawMetric("c419eb47-250e-44ec-98e1-f86b1a813520")
// 	utils.PanicError("TestGetRawMetrics", err)
// 	utils.Dump("TestGetRawMetrics", data)
// }

func TestCreateMint(t *testing.T) {
	var sign = &models.MintSign{
		Nonce:  2,
		Amount: "0xffaaa1",
		Iot:    iotAddr,
	}
	_, err := sign.Sign(testDomainMinter, iotPrv)
	utils.PanicError("TestCreateMint", err)

	utils.Dump("MintSign: ", sign)
}

func TestGetMint(t *testing.T) {
	var signeds, err = iotRepoTest.GetMintSigns(&domain.RIotGetMintSignList{
		From:  time.Now().Unix() - 60*86400,
		To:    time.Now().Unix(),
		IotId: 16,
	})
	utils.PanicError("TestGetMint", err)
	utils.Dump("TestGetMint", signeds)
}

func TestMint(t *testing.T) {
	var pk = "0123456789012345678901234567890123456789012345678901234567880000"
	var iotAddr = "0xe445517abb524002bb04c96f96abb87b8b19b53d"
	var amount1 = 9 * 1e9
	var amount2 = 12 * 1e9

	var sign1 = &models.MintSign{
		Nonce:  1,
		IotId:  292,
		Iot:    iotAddr,
		Amount: dmodels.NewBigNumber(int64(amount1)).ToHex(),
	}
	var sign2 = &models.MintSign{
		Nonce:  1,
		IotId:  292,
		Iot:    iotAddr,
		Amount: dmodels.NewBigNumber(int64(amount2)).ToHex(),
	}

	_, err := sign1.Sign(testDomainMinter, pk)
	utils.PanicError("", err)

	err = iotRepoTest.CreateMint(&domain.RIotMint{
		Nonce:  sign1.Nonce,
		Amount: sign1.Amount,
		Iot:    iotAddr,
		R:      sign1.R,
		S:      sign1.S,
		V:      sign1.V,
	})
	utils.PanicError("", err)

	time.Sleep(3 * time.Second)

	_, err = sign2.Sign(testDomainMinter, pk)
	utils.PanicError("", err)

	iotRepoTest.CreateMint(&domain.RIotMint{
		Iot:    iotAddr,
		Amount: sign2.Amount,
		Nonce:  sign2.Nonce,
		R:      sign2.R,
		S:      sign2.S,
		V:      sign2.V,
	})
	utils.PanicError("", err)
}

func TestGetMinted(t *testing.T) {
	var now = time.Now().Unix()
	log.Println(now-30*86400, now)
	data, err := iotRepoTest.GetMinted(&domain.RIotGetMintedList{
		From:     1693760400,
		To:       1693846799,
		IotId:    291,
		Interval: 1,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestIsIotActived(t *testing.T) {
	var now = time.Now().Unix()
	var to = now - 86400
	log.Println(to, now)
	data, err := iotRepoTest.IsIotActived(&domain.RIsIotActiced{
		From:  to,
		To:    now,
		IotId: 290,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}
