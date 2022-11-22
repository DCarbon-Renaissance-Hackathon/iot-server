package models

import (
	"time"

	"github.com/Dcarbon/iott-cloud/libs/esign"
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

type MintSign struct {
	ID        int64     `json:"id" gorm:"primary_key"` //
	Nonce     int64     `json:"nonce" gorm:"index"`    //
	Amount    string    `json:"amount" `               // Hex
	IOT       string    `json:"iot" gorm:"index"`      // IoT Address
	R         string    `json:"r" `                    //
	S         string    `json:"s" `                    //
	V         string    `json:"v" `                    //
	CreatedAt time.Time `json:"createdAt" `            //
	UpdatedAt time.Time `json:"updatedAt" `            //
}

func (*MintSign) TableName() string { return TableNameMintSign }

// prvStr: private key (hex)
func (msign *MintSign) Sign(prvStr string) ([]byte, error) {
	signedRaw, err := minterDomain.Sign(prvStr, map[string]interface{}{
		"iot":    msign.IOT,
		"amount": msign.Amount,
		"nonce":  msign.Nonce,
	})
	if nil != err {
		return nil, err
	}

	msign.R = hexutil.Encode(signedRaw[:32])
	msign.S = hexutil.Encode(signedRaw[32:64])
	msign.V = hexutil.Encode(signedRaw[64:])

	return signedRaw, nil
}

func (msign *MintSign) Verify() error {
	var data = map[string]interface{}{
		"iot":    msign.IOT,
		"amount": msign.Amount,
		"nonce":  msign.Nonce,
	}

	var signed, err = hexutil.Decode(
		esign.HexConcat(
			[]string{msign.R, msign.S, msign.V},
		),
	)

	if nil != err {
		return NewError(ECodeIOTInvalidMintSign, "Invalid mint sign: "+err.Error())
	}

	err = minterDomain.Verify(msign.IOT, signed, data)
	if nil != err {
		return NewError(ECodeIOTInvalidMintSign, "Invalid mint sign: "+err.Error())
	}
	return nil
}
