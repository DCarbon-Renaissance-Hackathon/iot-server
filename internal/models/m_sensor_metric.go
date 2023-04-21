package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var regString = regexp.MustCompile(`"*"$`)

type Sm struct {
	ID        string     `gorm:"primaryKey"`
	SignID    string     ``
	Indicator *AllMetric `gorm:"type:json"`
	CreatedAt time.Time  ``
}

// Sensor metric
type SmFloat struct {
	ID        string    ``
	SignID    string    ``
	Indicator float64   ``
	CreatedAt time.Time ``
}

func (*SmFloat) TableName() string { return TableNameSmFloat }

// Sensor metric gps
type SmGPS struct {
	ID        string     `json:"id"`
	SignID    string     `json:"signId"`
	Position  *Point4326 `json:"indicator" gorm:"type:geometry(POINT, 4326)"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (*SmGPS) TableName() string { return TableNameSmGPS }

// Sensor metric signature
type SmSignature struct {
	ID        string    `json:"id" `                  //
	IsIotSign bool      `json:"isIotSign" `           //
	IotID     int64     `json:"iotID" `               //
	SensorID  int64     `json:"sensorID" `            //
	Data      string    `json:"data" `                // Hex json of SensorMetricExtract
	Signed    string    `json:"signed" gorm:"unique"` // RSV Data
	CreatedAt time.Time `json:"createdAt" `           //
}

func (*SmSignature) TableName() string { return TableNameSmSignature }

func (sm *SmSignature) VerifySignature(addr EthAddress, sType SensorType) (*SMExtract, error) {
	err := addr.VerifyPersonalSign(sm.Data, sm.Signed)
	if nil != err {
		return nil, err
	}

	rawX, err := hexutil.Decode(sm.Data)
	if nil != err {
		return nil, NewError(ECodeInvalidSignature, "Data of signature must be hex")
	}

	x := &SMExtract{}
	err = json.Unmarshal(rawX, x)
	if nil != err {
		return nil, err
	}
	return x, x.IsValid(sType)
}

type SMExtract struct {
	From      int64      `json:"from"`
	To        int64      `json:"to"`
	Indicator AllMetric  `json:"indicator"`
	Address   EthAddress `json:"address"` // Sign (sensor or iot ) address
}

func (smx *SMExtract) IsValid(sType SensorType) error {
	if smx.From <= 1578104100 || smx.To > time.Now().Unix() {
		return NewError(ECodeSensorInvalidMetric, "Time range of metric is invalid [1578104100, now)")
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

type DefaultMetric struct {
	Val Float64 `json:"value,omitempty"`
}

type GPSMetric struct {
	Lat Float64 `json:"lat,omitempty"`
	Lng Float64 `json:"lng,omitempty"`
}

type AllMetric struct {
	DefaultMetric
	GPSMetric
}

func (am *AllMetric) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var rs = &AllMetric{}
	var err error
	switch vt := value.(type) {
	case string:
		if vt == `""` {
			return nil
		}
		err = json.Unmarshal([]byte(vt), rs)
	case []byte:
		err = json.Unmarshal(vt, rs)
	default:
		return errors.New("can't scan metric")
	}
	if nil != err {
		return err
	}
	if nil == am {
		am = new(AllMetric)
	}
	*am = *rs
	return nil
}

func (am AllMetric) Value() (driver.Value, error) {
	return json.Marshal(am)
}

func (am *AllMetric) IsValid(sType SensorType) error {
	switch sType {
	case SensorTypeFlow:
		if am.DefaultMetric.Val <= 0 {
			return NewError(ECodeSensorInvalidMetric, "Indicator of metric (value) must be > 0")
		}
	case SensorTypePower:
		if am.DefaultMetric.Val <= 0 {
			return NewError(ECodeSensorInvalidMetric, "Indicator of metric (value) must be > 0")
		}
	case SensorTypeGPS:
		if am.Lat == 0 && am.Lng == 0 {
			return NewError(ECodeSensorInvalidMetric, "Indicator of metric (gps) must be != 0")
		}
	}
	return nil
}

// type MetrictAggregate struct {
// 	ID        string    ``
// 	IotID     int64     ``
// 	Type      int64     ``
// 	Indicator float64   ``
// 	CreatedAt time.Time ``
// }

type Float64 float64

func (f *Float64) MarshalJSON() ([]byte, error) {
	if nil == f {
		return []byte("0"), nil
	}
	return []byte(fmt.Sprintf(`"%f"`, *f)), nil
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	var s = string(data)
	if regString.Match(data) {
		s = s[1 : len(s)-1]
	}

	v, err := strconv.ParseFloat(s, 64)
	if nil != err {
		return err
	}

	if nil == f {
		f = new(Float64)
	}

	*f = Float64(v)
	return nil
}
