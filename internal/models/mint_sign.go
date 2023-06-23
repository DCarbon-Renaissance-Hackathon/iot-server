package models

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// var minterDomain = esign.MustNewERC712(
// 	&esign.TypedDataDomain{
// 		Name:              "Carbon",
// 		Version:           "1",
// 		ChainId:           1,
// 		VerifyingContract: "0xA1E064Fd61B76cf11CE3b5816344f861b6318cea",
// 	},
// 	esign.MustNewTypedDataField(
// 		"Mint",
// 		esign.TypedDataStruct,
// 		esign.MustNewTypedDataField("iot", esign.TypedDataAddress),
// 		esign.MustNewTypedDataField("amount", "uint256"),
// 		esign.MustNewTypedDataField("nonce", "uint256"),
// 	),
// )

// const Precision = int64(1e9)

type MintSign struct {
	ID        int64     `json:"id" gorm:"primary_key"` //
	IotId     int64     `json:"iotId"`                 // IoT id
	Nonce     int64     `json:"nonce" gorm:"index"`    //
	Amount    string    `json:"amount" `               // Hex
	Iot       string    `json:"iot" gorm:"index"`      // IoT Address
	R         string    `json:"r" `                    //
	S         string    `json:"s" `                    //
	V         string    `json:"v" `                    //
	CreatedAt time.Time `json:"createdAt" `            //
	UpdatedAt time.Time `json:"updatedAt" `            //
}

func (*MintSign) TableName() string { return TableNameMintSign }

// Only for test
// pk: private key (hex)
func (msign *MintSign) Sign(dMinter *esign.ERC712, pk string) ([]byte, error) {
	signedRaw, err := dMinter.Sign(pk, map[string]interface{}{
		"iot":    msign.Iot,
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

func (msign *MintSign) Verify(dMinter *esign.ERC712) error {
	var data = map[string]interface{}{
		"iot":    msign.Iot,
		"amount": msign.Amount,
		"nonce":  msign.Nonce,
	}

	var signed, err = hexutil.Decode(
		esign.HexConcat(
			[]string{msign.R, msign.S, msign.V},
		),
	)

	if nil != err {
		return dmodels.NewError(dmodels.ECodeIOTInvalidMintSign, "Invalid mint sign: "+err.Error())
	}

	err = dMinter.Verify(msign.Iot, signed, data)
	if nil != err {
		return dmodels.NewError(dmodels.ECodeIOTInvalidMintSign, "Invalid mint sign: "+err.Error())
	}
	return nil
}

type Minted struct {
	ID        string    `json:"id,omitempty" `
	IotId     int64     `json:"iotId,omitempty" gorm:"index:minted_idx_ca_iot,priority:2"`
	Carbon    int64     `json:"carbon,omitempty" `
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"index:minted_idx_ca_iot,priority:1"`
}

func (*Minted) TableName() string { return TableNameMinted }
