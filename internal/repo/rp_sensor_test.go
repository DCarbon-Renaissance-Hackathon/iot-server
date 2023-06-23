package repo

import (
	"log"
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
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
		Status:  dmodels.DeviceStatusSuccess,
	},
	{
		ID:      2,
		Project: 2,
		Type:    models.IOTTypeFertilizer,
		Status:  dmodels.DeviceStatusSuccess,
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
	var addrs = []dmodels.EthAddress{
		dmodels.EthAddress("0xdC1A00c3cb7f769ED0C3021A38EC7cfCB5D0631e"),
		dmodels.EthAddress("0x973Fe93EcEA2F0A622377cC57FAb8EA596d18c63"),
		dmodels.EthAddress("0x69d1A0c44837bebA14b3F4dbb3384a546351E601"),
		dmodels.EthAddress("0xa45670F6d5bE173E07F911a435Dd83792E477D8F"),
		dmodels.EthAddress("0x0aedB9aCf69eB663BBAE23F1C8Eb8024da29fB71"),
		dmodels.EthAddress("0x451ea604180854155EAC73f82F1D36b80d648dE3"),
		dmodels.EthAddress("0x87A21119eb18DF1fFae01539D2B0AF6B39A508f2"),
		dmodels.EthAddress("0xc02050Ff0aF3E159E934067DA341e135441d60Fc"),
		dmodels.EthAddress("0x4bF02aF54d81BFBe2EdD68F1159B614A66b71201"),
		dmodels.EthAddress("0x07A88F400A4739F766B31833c8193621D4a8cc04"),
	}

	for _, it := range addrs {
		_, err := sensorImpl.CreateSensor(&domain.RCreateSensor{
			IotID:   iotTestSensors[0].ID,
			Type:    dmodels.SensorTypePower,
			Address: it,
		})
		utils.PanicError("", err)
	}
}

func TestSensorCreate2(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, err := sensorImpl.CreateSensor(&domain.RCreateSensor{
			IotID: iotTestSensors[1].ID,
			Type:  dmodels.SensorTypePower,
			// CreatedAt: time.Now(),
		})
		utils.PanicError("", err)
	}
}

func TestSensorChangeStatus(t *testing.T) {
	sensor, err := sensorImpl.ChangeSensorStatus(&domain.RChangeSensorStatus{
		ID:     31,
		Status: dmodels.DeviceStatusSuccess,
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
	var sensorAddr = dmodels.EthAddress("0x69d1a0c44837beba14b3f4dbb3384a546351e601")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567890002"
	var smx = &models.SMExtract{
		From: 1578104105,
		To:   1578104106,
		Indicator: &dmodels.AllMetric{
			DefaultMetric: dmodels.DefaultMetric{
				Val: 10.1,
			},
			GPSMetric: dmodels.GPSMetric{
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
	var iotAddr = dmodels.EthAddress("0xE445517AbB524002Bb04C96F96aBb87b8B19b53d")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880000"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		Indicator: &dmodels.AllMetric{
			DefaultMetric: dmodels.DefaultMetric{
				Val: 10.1,
			},
		},
		Address: iotAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)

	data, err := sensorImpl.CreateSensorMetric(&domain.RCreateSensorMetric{
		Data:        signed.Data,
		Signed:      signed.Signed,
		SignAddress: iotAddr,
		IsIotSign:   true,
		SensorID:    31,
		IotID:       2,
	})
	utils.PanicError("CreateSM", err)
	utils.Dump("SM", data)
}

func TestSensorGetSM(t *testing.T) {
	var now = time.Now().Unix()
	data, err := sensorImpl.GetMetrics(&domain.RGetSM{
		From:  now - 1000,
		To:    now,
		IotId: 277,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestGenerateSignMetric(t *testing.T) {
	var iotAddr = dmodels.EthAddress("0xE445517AbB524002Bb04C96F96aBb87b8B19b53d")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880000"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		Indicator: &dmodels.AllMetric{
			DefaultMetric: dmodels.DefaultMetric{
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
	var iotAddr = dmodels.EthAddress("0x19Adf96848504a06383b47aAA9BbBC6638E81afD")
	var pKey = "0x0123456789012345678901234567890123456789012345678901234567880001"
	var smx = &models.SMExtract{
		From: 1578104103,
		To:   1578104104,
		// Indicator: 10.1,
		Address: iotAddr,
	}
	var signed, err = smx.Signed(pKey)
	utils.PanicError("", err)

	var req = &domain.RCreateSensorMetric{
		Data:        signed.Data,
		Signed:      signed.Signed,
		SignAddress: iotAddr,
		SensorID:    31,
		IsIotSign:   true,
		IotID:       1,
	}
	utils.Dump("Request", req)
}

func TestGetExtract(t *testing.T) {
	var sign = &models.SmSignature{
		ID:     "4759d9b0-c948-4533-bd5a-128d2bc0bcdc",
		Data:   "0x7b2266726f6d223a313638383131383036382c22746f223a313638383131383036392c22696e64696361746f72223a7b2276616c7565223a22393033343132222c226c6174223a2232312e353033373133222c226c6e67223a224e614e227d2c2261646472657373223a22307837326566396461326166316436353762336664313665393366623965366438326334633631356630227d",
		Signed: "0x7c3732463c11608ce59dc7f50e444a833c308a626e3528ec4efbd31ce3b59e0866285294bcecd1e1b3734ca247cdb5f0f73958d4a34e627a48b10a692823ead31c",
	}
	data, err := sign.ExtractData()
	utils.PanicError("", err)
	utils.Dump("", data)
}

func TestMigrate(t *testing.T) {
	var interval = 3600
	var nowUnix = time.Now().Unix()
	var from = time.Now().Add(-1 * 250 * 86400 * time.Second).Unix()
	var to = from + int64(interval)

	for {
		log.Printf("From %+v  && to: %+v\n", time.Unix(from, 0), time.Unix(to, 0))
		signedMetrics, err := sensorImpl.GetSignedMetric(&domain.RGetSM{
			From: from,
			To:   to,
		})
		utils.PanicError("Get signed metric ", err)

		for _, it := range signedMetrics {
			// log.Println(it)
			_, err = sensorImpl.migrateSM(it)
			utils.PanicError("Migrate sm", err)
			// _, err := it.ExtractData()
			// if nil != err {
			// 	log.Println("Migrate error: ", it.ID, err)
			// }
		}
		if from > nowUnix {
			break
		}

		from = to + 1
		to = from + int64(interval)
	}
}

func TestMetricAggregate(t *testing.T) {
	var nowUnix = time.Now().Unix()
	var from = nowUnix - 3*86400
	var to = nowUnix
	log.Println(from, to)

	var data, err = sensorImpl.getMetricAggregate(&domain.RSMAggregate{
		From:     from,
		To:       to,
		IotId:    291,
		SensorId: 76,
		Interval: 2,
	})
	utils.PanicError("", err)
	utils.Dump("", data)
}
