package esign

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const messagePrefix = "\x19Ethereum Signed Message:\n"

func VerifyPersonalSign(address string, org, signed []byte) error {
	var raw = []byte(fmt.Sprintf(messagePrefix+"%d", len(org)))
	raw = append(raw, org...)
	return Verify(address, raw, signed)
}

func SignPersonal(prvStr string, data []byte) ([]byte, error) {
	var raw = []byte(fmt.Sprintf(messagePrefix+"%d", len(data)))
	raw = append(raw, data...)
	return Sign(prvStr, raw)
}

func Sign(prvStr string, data []byte) ([]byte, error) {
	prv, err := crypto.HexToECDSA(prvStr)
	if nil != err {
		return nil, err
	}

	hash := crypto.Keccak256Hash(data)
	signed, err := crypto.Sign(hash[:], prv)
	if nil != err {
		return nil, err
	}
	signed[64] += 27
	return signed, nil
}

func Verify(addrStr string, org, signed []byte) error {
	inputPub := common.HexToAddress(addrStr)

	hash := crypto.Keccak256Hash(org)
	if len(signed) != 65 || (signed[64] != 27 && signed[64] != 28) {
		return errors.New("invalid signature")
	}
	var signed2 = make([]byte, len(signed))
	copy(signed2, signed)

	signed2[64] -= 27
	sigPubKey, err := crypto.Ecrecover(hash[:], signed2)
	if nil != err {
		return err
	}

	pubkey, err := crypto.UnmarshalPubkey(sigPubKey)
	if nil != err {
		return err
	}
	// fmt.Println("Address: ", crypto.PubkeyToAddress(*pubkey))
	if crypto.PubkeyToAddress(*pubkey) == inputPub {
		return nil
	}
	return errors.New("not match")
}

func GenerateKey() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)
}
