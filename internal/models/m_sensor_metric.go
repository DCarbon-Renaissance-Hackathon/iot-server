package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/ecodes"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// var regString = regexp.MustCompile(`"*"$`)

// Sensor metric data
type Sm struct {
	ID        string             `json:"id"       gorm:"primaryKey"`          //
	SignID    string             `json:"signId"`                              //
	SensorID  int64              `json:"sensorID" gorm:"index_ca,priority:3"` //
	IotID     int64              `json:"iotID"    gorm:"index_ca,priority:2"` //
	Indicator *dmodels.AllMetric `json:"metric"   gorm:"type:json"`           //
	CreatedAt time.Time          `json:"ca"       gorm:"index_ca,priority:1"` //
}

func (*Sm) TableName() string { return TableNameSm }

// Sensor metric
// type SmFloat struct {
// 	ID        string    ``
// 	SignID    string    ``
// 	Indicator float64   ``
// 	CreatedAt time.Time ``
// }

// func (*SmFloat) TableName() string { return TableNameSmFloat }

// // Sensor metric gps
// type SmGPS struct {
// 	ID        string     `json:"id"`
// 	SignID    string     `json:"signId"`
// 	Position  *Point4326 `json:"indicator" gorm:"type:geometry(POINT, 4326)"`
// 	CreatedAt time.Time  `json:"createdAt"`
// }

// func (*SmGPS) TableName() string { return TableNameSmGPS }

// Sensor metric signature
type SmSignature struct {
	ID        string    `json:"id" `                                        //
	IsIotSign bool      `json:"isIotSign" `                                 //
	IotID     int64     `json:"iotID" gorm:"index:sms_index_ca,priority:2"` //
	SensorID  int64     `json:"sensorID" `                                  //
	Data      string    `json:"data" `                                      // Hex json of SensorMetricExtract
	Signed    string    `json:"signed" gorm:"unique"`                       // RSV Data
	CreatedAt time.Time `json:"ca" gorm:"index:sms_index_ca,priority:1"`    //
}

func (*SmSignature) TableName() string { return TableNameSmSignature }

func (sm *SmSignature) VerifySignature(addr dmodels.EthAddress, sType dmodels.SensorType) (*SMExtract, error) {
	err := addr.VerifyPersonalSign(sm.Data, sm.Signed)
	if nil != err {
		return nil, err
	}

	x, err := sm.ExtractData()
	if nil != err {
		return nil, err
	}

	return x, x.IsValid(sType)
}

func (sm *SmSignature) ExtractData() (*SMExtract, error) {
	rawX, err := hexutil.Decode(sm.Data)
	if nil != err {
		return nil, dmodels.NewError(ecodes.InvalidSignature, "Data of signature must be hex")
	}

	x := &SMExtract{}
	err = json.Unmarshal(rawX, x)
	if nil != err {
		log.Println("Extract data error: ", err)
		log.Println("data: ", string(rawX))
		return nil, err
	}
	return x, nil

}

type SMExtract struct {
	From      int64              `json:"from"`
	To        int64              `json:"to"`
	Indicator *dmodels.AllMetric `json:"indicator"`
	Address   dmodels.EthAddress `json:"address"` // Sign (sensor or iot ) address
}

func (smx *SMExtract) IsValid(sType dmodels.SensorType) error {
	if smx.From <= 1578104100 || smx.To > time.Now().Unix() {
		return dmodels.NewError(ecodes.SensorInvalidMetric, "Time range of metric is invalid [1578104100, now)")
	}

	err := smx.Indicator.IsValid(sType)
	if nil != err {
		return err
	}

	// if smx.Indicator <= 0 {
	// 	return NewError(ECodeSensorInvalidMetric, "Indicator of metric must be gt 0")
	// }

	return nil
}

func (smx *SMExtract) Signed(pkey string) (*SmSignature, error) {
	raw, err := json.Marshal(smx)
	if nil != err {
		return nil, err
	}

	signedRaw, err := esign.SignPersonal(pkey, raw)
	if nil != err {
		return nil, err
	}

	return &SmSignature{
		IsIotSign: true,
		Data:      hexutil.Encode(raw),
		Signed:    hexutil.Encode(signedRaw),
	}, nil
}

// Instant sensor metric extract
type ISMExtract struct {
	Signer    dmodels.EthAddress `json:"signer"`    // Sign address (sensor or iot )
	CreatedAt int64              `json:"ca"`        //
	Indicator *dmodels.AllMetric `json:"indicator"` //
}

// Accumulate sensor metric
type ASMExtract struct {
	Signer dmodels.EthAddress `json:"signer"` // Sign address (sensor or iot )
	From   int64              `json:"from"`
	To     int64              `json:"to"`
	Value  dmodels.Float64    `json:"value"`
}

// type DefaultMetric struct {
// 	Val Float64 `json:"value,omitempty"`
// }

// type GPSMetric struct {
// 	Lat Float64 `json:"lat,omitempty"`
// 	Lng Float64 `json:"lng,omitempty"`
// }

// type AllMetric struct {
// 	DefaultMetric
// 	GPSMetric
// } // @name models.AllMetric

// func (am *AllMetric) Scan(value interface{}) error {
// 	if value == nil {
// 		return nil
// 	}

// 	var rs = &AllMetric{}
// 	var err error
// 	switch vt := value.(type) {
// 	case string:
// 		if vt == `""` {
// 			return nil
// 		}
// 		err = json.Unmarshal([]byte(vt), rs)
// 	case []byte:
// 		err = json.Unmarshal(vt, rs)
// 	default:
// 		return errors.New("can't scan metric")
// 	}
// 	if nil != err {
// 		return err
// 	}
// 	if nil == am {
// 		am = new(AllMetric)
// 	}
// 	*am = *rs
// 	return nil
// }

// func (am AllMetric) Value() (driver.Value, error) {
// 	return json.Marshal(am)
// }

// func (am *AllMetric) IsValid(sType dmodels.SensorType) error {
// 	switch sType {
// 	case dmodels.SensorTypeFlow:
// 		if am.DefaultMetric.Val <= 0 {
// 			return dmodels.NewError(dmodels.ECodeSensorInvalidMetric, "Indicator of metric (value) must be > 0")
// 		}
// 	case dmodels.SensorTypePower:
// 		if am.DefaultMetric.Val <= 0 {
// 			return dmodels.NewError(dmodels.ECodeSensorInvalidMetric, "Indicator of metric (value) must be > 0")
// 		}
// 	case dmodels.SensorTypeGPS:
// 		if am.Lat == 0 && am.Lng == 0 {
// 			return dmodels.NewError(dmodels.ECodeSensorInvalidMetric, "Indicator of metric (gps) must be != 0")
// 		}
// 	}
// 	return nil
// }

// type Float64 float64

// func (f *Float64) MarshalJSON() ([]byte, error) {
// 	if nil == f {
// 		return []byte("0"), nil
// 	}
// 	return []byte(fmt.Sprintf(`"%f"`, *f)), nil
// }

// func (f *Float64) UnmarshalJSON(data []byte) error {
// 	var s = string(data)
// 	if regString.Match(data) {
// 		s = s[1 : len(s)-1]
// 	}

// 	v, err := strconv.ParseFloat(s, 64)
// 	if nil != err {
// 		return err
// 	}

// 	if nil == f {
// 		f = new(Float64)
// 	}

// 	*f = Float64(v)
// 	return nil
// }
