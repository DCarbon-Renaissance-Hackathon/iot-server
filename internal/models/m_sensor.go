package models

import (
	"encoding/json"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type SensorType int32

const (
	SensorTypeNone  SensorType = 0
	SensorTypeFlow  SensorType = 1
	SensorTypePower SensorType = 2
)

type SensorStatus int32

const (
	SensorStatusReject   SensorStatus = -1
	SensorStatusRegister SensorStatus = 0
	SensorStatusSuccess  SensorStatus = 10
)

type Sensor struct {
	ID        int64        ``
	IotID     int64        ``
	Address   *EthAddress  `gorm:"index:,unique,where:length(address) > 0"`
	Type      SensorType   `` // CH4, KW, MW, ...
	Status    SensorStatus ``
	CreatedAt time.Time    ``
}

func (*Sensor) TableName() string { return TableNameSensors }

// Sensor metric
type SM struct {
	ID        string    ``
	SignID    string    ``
	Indicator float64   ``
	CreatedAt time.Time ``
}

func (*SM) TableName() string { return TableNameSM }

// Sensor metric signature
type SMSignature struct {
	ID        string    ``
	IsIotSign bool      ``
	IotID     int64     ``
	SensorID  int64     ``
	Data      string    `` // Hex json of SensorMetricExtract
	Signed    string    `` // RSV Data
	CreatedAt time.Time `` //
}

func (*SMSignature) TableName() string { return TableNameSMSignature }

func (sm *SMSignature) VerifySignature(addr EthAddress) (*SMExtract, error) {
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

	return x, x.IsValid()
}

type SMExtract struct {
	From      int64      `json:"from"`
	To        int64      `json:"to"`
	Indicator float64    `json:"indicator"`
	Address   EthAddress `json:"address"` // Sign (sensor or iot ) address
}

func (smx *SMExtract) IsValid() error {
	if smx.From <= 1578104100 || smx.To > time.Now().Unix() {
		return NewError(ECodeSensorInvalidMetric, "Time range of metric is invalid")
	}

	if smx.Indicator <= 0 {
		return NewError(ECodeSensorInvalidMetric, "Indicator of metric must be gt 0")
	}

	return nil
}

func (smx *SMExtract) Signed(pkey string) (*SMSignature, error) {
	raw, err := json.Marshal(smx)
	if nil != err {
		return nil, err
	}

	signedRaw, err := esign.SignPersonal(pkey, raw)
	if nil != err {
		return nil, err
	}

	return &SMSignature{
		IsIotSign: true,
		Data:      hexutil.Encode(raw),
		Signed:    hexutil.Encode(signedRaw),
	}, nil
}

// type MetrictAggregate struct {
// 	ID        string    ``
// 	IotID     int64     ``
// 	Type      int64     ``
// 	Indicator float64   ``
// 	CreatedAt time.Time ``
// }
