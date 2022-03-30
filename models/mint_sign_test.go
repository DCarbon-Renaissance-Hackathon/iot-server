package models

import (
	"log"
	"testing"

	"github.com/Dcarbon/libs/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const AddrStr = "0xCC719739eD48B0258456F104DA7ba83Ba6881C35"
const PrvStr = "5763b65df1b1860bfa8a372ae589f1a67811c3e4a7234d29fc3d68d2c531e547"

func TestMintSignAndVerify(t *testing.T) {
	var m = &MintSign{
		ID:     0,
		Nonce:  100,
		Amount: "0xaabbcc",
		IOT:    AddrStr,
	}
	signed, err := m.Sign(PrvStr)
	utils.PanicError("TestMintSign", err)

	log.Println("Sign: ", hexutil.Encode(signed))

	m.R = hexutil.Encode(signed[:32])
	m.S = hexutil.Encode(signed[32:64])
	m.V = hexutil.Encode(signed[64:])

	err = m.Verify()
	utils.PanicError("TestMintVerify", err)
}
