package models

import (
	"encoding/json"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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
	ID        string    `json:"id" `                  //
	IsIotSign bool      `json:"isIotSign" `           //
	IotID     int64     `json:"iotID" `               //
	SensorID  int64     `json:"sensorID" `            //
	Data      string    `json:"data" `                // Hex json of SensorMetricExtract
	Signed    string    `json:"signed" gorm:"unique"` // RSV Data
	CreatedAt time.Time `json:"createdAt" `           //
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
