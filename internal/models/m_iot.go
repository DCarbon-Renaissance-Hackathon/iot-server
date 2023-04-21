package models

type IOTStatus int

const (
	IOTStatusReject   IOTStatus = -1
	IOTStatusRegister IOTStatus = 0
	IOTStatusSuccess  IOTStatus = 10
)

type IOTType int

const (
	IOTTypeNone        IOTType = 0
	IOTTypeWindPower   IOTType = 10
	IOTTypeSolarPower  IOTType = 11
	IOTTypeBurnMethane IOTType = 20
	IOTTypeFertilizer  IOTType = 30
	IOTTypeTrash       IOTType = 31
)

type IOTDevice struct {
	ID       int64      `json:"id" gorm:"primary_key"`
	Project  int64      `json:"project" `
	Address  EthAddress `json:"address" gorm:"unique"`
	Type     IOTType    `json:"type" `
	Status   IOTStatus  `json:"status"`
	Position Point4326  `json:"position" gorm:"type:geometry(POINT, 4326)"`
}

func (*IOTDevice) TableName() string { return TableNameIOT }

// type ExtractMetric struct {
// 	ID       string            ``
// 	IsResult bool              ``
// 	Warning  int               `` // Warning code
// 	From     int64             `json:"from"`
// 	To       int64             `json:"to"`
// 	Position Point4326         `json:"pos" gorm:"column:pos;index;type:geometry(POINT, 4326)"`
// 	Metrics  dbutils.MapSFloat `json:"metrics" gorm:"type:json"` // Unit m3/s Ex: {"ch4": 1.1}
// }

// type Metric struct {
// 	ID        string        `json:"id,omitempty"`                      //
// 	Address   string        `json:"address,omitempty"`                 // IOT address
// 	Data      string        `json:"data,omitempty"`                    // Json string
// 	Signed    string        `json:"signed,omitempty"`                  // Hex
// 	Extract   ExtractMetric `json:"extract,omitempty" gorm:"embedded"` //
// 	CreatedAt time.Time     `json:"createdAt,omitempty"`               //
// }

// func (*Metric) TableName() string { return TableNameMetrics }

// func (m *Metric) VerifySignature() error {
// 	rawOrg, err := hexutil.Decode(m.Data)
// 	if nil != err {
// 		return NewError(ECodeInvalidSignature, "Data of signature must be hex")
// 	}

// 	rawSigned, err := hexutil.Decode(m.Signed)
// 	if nil != err {
// 		return NewError(ECodeInvalidSignature, "Signature must be hex")
// 	}

// 	err = esign.VerifyPersonalSign(m.Address, rawOrg, rawSigned)
// 	if nil != err {
// 		return NewError(ECodeInvalidSignature, "Signature invalid")
// 	}
// 	return nil
// }
