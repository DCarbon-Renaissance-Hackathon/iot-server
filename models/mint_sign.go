package models

import (
	"time"

	"github.com/Dcarbon/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var minterDomain = esign.MustNewERC712(
	&esign.TypedDataDomain{
		Name:              "Carbon",
		Version:           "1",
		ChainId:           1,
		VerifyingContract: "0xA1E064Fd61B76cf11CE3b5816344f861b6318cea",
	},
	esign.MustNewTypedDataField(
		"Mint",
		esign.TypedDataStruct,
		esign.MustNewTypedDataField("iot", esign.TypedDataAddress),
		esign.MustNewTypedDataField("amount", "uint256"),
		esign.MustNewTypedDataField("nonce", "uint256"),
	),
)

type MintSignature struct {
	ID        int64     `gorm:"primary_key"` //
	Nonce     int64     `gorm:"index"`       //
	Amount    string    ``                   // Hex
	IoT       string    `gorm:"index"`       // IoT Address
	R         string    ``                   //
	S         string    ``                   //
	V         string    ``                   //
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}

func (*MintSignature) TableName() string { return TableNameMintSignature }

// prvStr: private key (hex)
func (msign *MintSignature) Sign(prvStr string) ([]byte, error) {
	return minterDomain.Sign(prvStr, map[string]interface{}{
		"iot":    msign.IoT,
		"amount": msign.Amount,
		"nonce":  msign.Nonce,
	})
}

func (msign *MintSignature) Verify() error {
	var data = map[string]interface{}{
		"iot":    msign.IoT,
		"amount": msign.Amount,
		"nonce":  msign.Nonce,
	}
	var signed, err = hexutil.Decode(
		esign.HexConcat(
			[]string{msign.R, msign.S, msign.V},
		),
	)
	if nil != err {
		return err
	}
	return minterDomain.Verify(msign.IoT, signed, data)
}
