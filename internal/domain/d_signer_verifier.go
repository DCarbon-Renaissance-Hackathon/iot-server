package domain

import (
	"fmt"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type SignedToken struct {
	Address  dmodels.EthAddress `json:"address"`  // Sign address
	SignedAt int64              `json:"signedAt"` // Timestamp (second)
	Signed   string             `json:"signed"`   // Hex string
}

// Verify token by address
type ISignerVerifier interface {
	IsValid(token *SignedToken) error
}

type verifier struct {
}

func NewVerifier() ISignerVerifier {
	var v = &verifier{}
	return v
}

func (v *verifier) IsValid(token *SignedToken) error {
	if token.SignedAt > time.Now().Unix() {
		return dmodels.ErrBadRequest("Signature too early ")
	}

	if token.SignedAt < time.Now().Unix()-4300 {
		return dmodels.ErrBadRequest("Signature was expired")
	}

	var org = fmt.Sprintf("dcarbon_%d_%s", token.SignedAt, token.Address)
	var signedBytes, err = hexutil.Decode(token.Signed)
	if nil != err {
		return dmodels.ErrBadRequest("Invalid sign " + err.Error())
	}

	err = esign.VerifyPersonalSign(string(token.Address), []byte(org), signedBytes)
	if nil != err {
		return dmodels.ErrBadRequest("Invalid signed" + err.Error())
	}
	return nil
}
