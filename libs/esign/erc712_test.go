package esign

import (
	"log"
	"math/big"
	"testing"
)

func TestErc712(t *testing.T) {
	var minter, err = NewERC712(
		&TypedDataDomain{
			Name:              "Carbon",
			Version:           "1",
			ChainId:           1,
			VerifyingContract: "0xA1E064Fd61B76cf11CE3b5816344f861b6318cea",
		},
		MustNewTypedDataField(
			"Mint",
			TypedDataStruct,
			MustNewTypedDataField("iot", TypedDataAddress),
			MustNewTypedDataField("amount", "uint256"),
			MustNewTypedDataField("nonce", "uint256"),
		),
	)
	panicError("", err)

	log.Println("DomainHash: ", minter.domainHash)
	hash, err := minter.Hash(map[string]interface{}{
		"iot":    "0x5c77E37aA7AFa0064b1eFb01cFbf2EfdFF49E7EA",
		"amount": 101,
		"nonce":  72727269,
	})
	panicError("Minter hash", err)
	log.Println("Minter hash: ", hash)
}

func TestInt(t *testing.T) {
	var typeInt = MustNewTypedDataField("test_int", "uint256")

	hash, err := typeInt.Encode(big.NewInt(101))
	panicError("", err)
	log.Println(hash)
}

func panicError(label string, err error) {
	if nil != err {
		panic(label + " error: " + err.Error())
	}
}
