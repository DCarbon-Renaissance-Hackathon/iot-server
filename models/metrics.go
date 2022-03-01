package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/Dcarbon/libs/dbutils"
	"github.com/Dcarbon/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TrackingType int

const (
	TrackingTypeFlow TrackingType = 1
)

type ExtractMetric struct {
	From     int64             `json:"from"`
	To       int64             `json:"to"`
	Position Point4326         `json:"pos" gorm:"column:pos;index;type:geometry(POINT, 4326)"`
	Metrics  dbutils.MapSFloat `json:"metrics" gorm:"type:json"` // Unit m3/s Ex: {"ch4": 1.1}
}

func (m *ExtractMetric) Scan(value interface{}) error {
	if nil == m {
		m = new(ExtractMetric)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSFloat invalid")
}

func (m ExtractMetric) Value() (driver.Value, error) {
	return json.Marshal(&m)
}

type Metric struct {
	ID        string        `json:"id,omitempty"`                       //
	Address   string        `json:"address,omitempty"`                  //
	Data      string        `json:"data,omitempty"`                     // Hex
	Signed    string        `json:"signed,omitempty"`                   // Hex
	Extract   ExtractMetric `json:"extract,omitempty" gorm:"type:json"` //
	CreatedAt time.Time     `json:"createdAt,omitempty"`                //
}

func (*Metric) TableName() string { return TableNameMetrics }

func (m *Metric) VerifySignature() error {
	rawOrg, err := hexutil.Decode(m.Data)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Data of signature must be hex")
	}

	rawSigned, err := hexutil.Decode(m.Signed)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Signature must be hex")
	}

	err = esign.VerifyPersonalSign(m.Address, rawOrg, rawSigned)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Signature invalid")
	}
	return nil
}
