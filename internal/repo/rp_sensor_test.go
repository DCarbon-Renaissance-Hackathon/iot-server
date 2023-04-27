package repo

import (
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

var sensorImpl *SensorRepo

var iotTestSensors = []*models.IOTDevice{
	{
		ID:      1,
		Project: 1,
		Type:    models.IOTTypeBurnMethane,
		Address: "0xE445517AbB524002Bb04C96F96aBb87b8B19b53d",
		Status:  models.DeviceStatusSuccess,
	},
	{
		ID:      2,
		Project: 2,
		Type:    models.IOTTypeFertilizer,
		Status:  models.DeviceStatusSuccess,
		Address: "0x19Adf96848504a06383b47aAA9BbBC6638E81afD",
	},
}

func init() {
	var err error

	// err := rss.InitResource(dbUrlTest, redisUrl)
	// utils.PanicError("", err)

	sensorImpl, err = NewSensorRepo()
	utils.PanicError("", err)
}

func TestSensorCreate(t *testing.T) {
	var addrs = []models.EthAddress{
		models.EthAddress("0xdC1A00c3cb7f769ED0C3021A38EC7cfCB5D0631e"),
		models.EthAddress("0x973Fe93EcEA2F0A622377cC57FAb8EA596d18c63"),
		models.EthAddress("0x69d1A0c44837bebA14b3F4dbb3384a546351E601"),
		models.EthAddress("0xa45670F6d5bE173E07F911a435Dd83792E477D8F"),
		models.EthAddress("0x0aedB9aCf69eB663BBAE23F1C8Eb8024da29fB71"),
		models.EthAddress("0x451ea604180854155EAC73f82F1D36b80d648dE3"),
		models.EthAddress("0x87A21119eb18DF1fFae01539D2B0AF6B39A508f2"),
		models.EthAddress("0xc02050Ff0aF3E159E934067DA341e135441d60Fc"),
		models.EthAddress("0x4bF02aF54d81BFBe2EdD68F1159B614A66b71201"),
		models.EthAddress("0x07A88F400A4739F766B31833c8193621D4a8cc04"),
	}

	for _, it := range addrs {
		_, err := sensorImpl.CreateSensor(&domain.RCreateSensor{
			IotID:   iotTestSensors[0].ID,
			Type:    models.SensorTypePower,
			Address: it,
		})
		utils.PanicError("", err)
	}
}

func TestSensorCreate2(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, err := sensorImpl.CreateSensor(&domain.RCreateSensor{
			IotID: iotTestSensors[1].ID,
			Type:  models.SensorTypePower,
			// CreatedAt: time.Now(),
		})
		utils.PanicError("", err)
	}
}

func TestSensorChangeStatus(t *testing.T) {
	sensor, err := sensorImpl.ChangeSensorStatus(&domain.RChangeSensorStatus{
		ID:     31,
		Status: models.DeviceStatusSuccess,
	})
	utils.PanicError("", err)
	utils.Dump("Changed sensor", sensor)
}

func TestSensorGetSensor(t *testing.T) {
	sensor, err := sensorImpl.GetSensor(&domain.SensorID{
		// ID:      1,
		Address: "0xdC1A00c3cb7f769ED0C3021A38EC7cfCB5D0631e",
	})
	utils.PanicError("", err)
	utils.Dump("Changed sensor", sensor)
}

func TestSensorGetSensors(t *testing.T) {
	data, err := sensorImpl.GetSensors(&domain.RGetSensors{
		Skip:  0,
		Limit: 3,
	})
	utils.PanicError("", err)
	utils.Dump("Changed sensor", data)
}

func TestSensorCreateSM(t *testing.T) {
	var sensorAddr = models.EthAddress("0x69d1a0c44837beba14b3f4dbb3384a546351e601")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567890002"
	var smx = &models.SMExtract{
		From: 1578104105,
		To:   1578104106,
		Indicator: &models.AllMetric{
			DefaultMetric: models.DefaultMetric{
				Val: 10.1,
			},
			GPSMetric: models.GPSMetric{
				Lng: 105.1,
				Lat: 22.1,
			},
		},
		Address: sensorAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)

	data, err := sensorImpl.CreateSM(&domain.RCreateSM{
		SensorAddress: sensorAddr,
		Data:          signed.Data,
		Signed:        signed.Signed,
	})
	utils.PanicError("CreateSM", err)
	utils.Dump("SM", data)

}

func TestSensorCreateSMFromIOT(t *testing.T) {
	var iotAddr = models.EthAddress("0xE445517AbB524002Bb04C96F96aBb87b8B19b53d")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880000"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		Indicator: &models.AllMetric{
			DefaultMetric: models.DefaultMetric{
				Val: 10.1,
			},
		},
		Address: iotAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)

	data, err := sensorImpl.CreateSMFromIot(&domain.RCreateSMFromIOT{
		SensorID:   31,
		IotID:      2,
		IotAddress: iotAddr,
		Data:       signed.Data,
		Signed:     signed.Signed,
	})
	utils.PanicError("CreateSM", err)
	utils.Dump("SM", data)
}

func TestSensorGetSM(t *testing.T) {
	data, err := sensorImpl.GetMetrics(&domain.RGetSM{
		From:  1578104100,
		To:    time.Now().Unix(),
		IotId: 1,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestGenerateSignMetric(t *testing.T) {
	var iotAddr = models.EthAddress("0xE445517AbB524002Bb04C96F96aBb87b8B19b53d")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880000"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		Indicator: &models.AllMetric{
			DefaultMetric: models.DefaultMetric{
				Val: 10.1,
			},
		},
		Address: iotAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)
	utils.Dump("signed: ", signed)
	utils.Dump("", smx)
}

func TestGenerateSignMetric2(t *testing.T) {
	// var iotID = 2
	var iotAddr = models.EthAddress("0x19Adf96848504a06383b47aAA9BbBC6638E81afD")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880001"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		// Indicator: 10.1,
		Address: iotAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)

	var req = &domain.RCreateSMFromIOT{
		Data:       signed.Data,
		Signed:     signed.Signed,
		IotAddress: iotAddr,
		SensorID:   31,
	}
	utils.Dump("Request", req)
}

func TestXxx(t *testing.T) {
	utils.Dump("", &models.AllMetric{})
}
