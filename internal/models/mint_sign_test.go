package models

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const AddrStr = "0x19Adf96848504a06383b47aAA9BbBC6638E81afD"
const PrvStr = "0123456789012345678901234567890123456789012345678901234567880001"

var testDomainMinter = esign.MustNewERC712(
	&esign.TypedDataDomain{
		Name:              "CARBON",
		Version:           "1",
		ChainId:           1337,
		VerifyingContract: "0x7BDDCb9699a3823b8B27158BEBaBDE6431152a85",
	},
	esign.MustNewTypedDataField(
		"Mint",
		esign.TypedDataStruct,
		esign.MustNewTypedDataField("iot", esign.TypedDataAddress),
		esign.MustNewTypedDataField("amount", "uint256"),
		esign.MustNewTypedDataField("nonce", "uint256"),
	),
)

func TestMintSignAndVerify(t *testing.T) {
	var m = &MintSign{
		ID:     0,
		Nonce:  3,
		Amount: "0xaabbccddee",
		IOT:    AddrStr,
	}
	signed, err := m.Sign(testDomainMinter, PrvStr)
	utils.PanicError("TestMintSign", err)

	log.Println("Sign: ", hexutil.Encode(signed))

	m.R = hexutil.Encode(signed[:32])
	m.S = hexutil.Encode(signed[32:64])
	m.V = hexutil.Encode(signed[64:])

	err = m.Verify(testDomainMinter)
	utils.PanicError("TestMintVerify", err)
	utils.Dump("Mint sign", m)
}

func TestFloat(t *testing.T) {
	var d = &DefaultMetric{
		Val: 100.1,
	}

	raw, err := json.Marshal(d)
	utils.PanicError("", err)

	var d2 = &DefaultMetric{}
	err = json.Unmarshal(raw, d2)
	utils.PanicError("", err)
}

func TestMintVerify(t *testing.T) {
	var m = &MintSign{
		ID:     0,
		Nonce:  3,
		Amount: "0xaabbccddee",
		R:      "0xca5979c9d43300870fa6c513e36ba4617179ed1b58ab7fcc8c6a19074b8c09ad",
		S:      "0x6d2ccc2c822c54c41c15915590c6f107f8c8ed938557e509e01bd121865a8acb",
		V:      "0x1c",
		IOT:    "0x19Adf96848504a06383b47aAA9BbBC6638E81afD",
	}

	err := m.Verify(testDomainMinter)
	utils.PanicError("TestMintVerify", err)
}
