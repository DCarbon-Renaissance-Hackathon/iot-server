package esign

import (
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var data = []byte("this is test")

const AddrStr = "0xCC719739eD48B0258456F104DA7ba83Ba6881C35"
const PrvStr = "5763b65df1b1860bfa8a372ae589f1a67811c3e4a7234d29fc3d68d2c531e547"

func TestSign(t *testing.T) {
	signed, err := Sign(PrvStr, data)
	if nil != err {
		panic(err)
	}
	log.Println("Signed: ", hexutil.Encode(signed))
	err = Verify(AddrStr, data, signed)
	if nil != err {
		panic(err)
	}

	// var rawSigned = "0x8386b2ea54797aa531d9b126f0322904852737dd59b798f5564f4b4ccc2312ac462b3af45fbd153e7d7cbb3a75226f45373f4f0e1c625a1136eb67c506b72ff81b"
	// signed, err := hexutil.Decode(rawSigned)
	// if nil != err {
	// 	panic(err)
	// }
	// err = VerifyPersonalSign(PubStr, data, signed)
	// if nil != err {
	// 	panic(err)
	// }
	// log.Println("Verify success")
}

func TestPersonal(t *testing.T) {
	var data = []byte("this is test")
	var signed, err = SignPersonal(PrvStr, data)
	if nil != err {
		panic("Sign personal error: " + err.Error())
	}
	log.Println("Signed: ", hexutil.Encode(signed))
	rd, _ := hexutil.Decode(hexutil.Encode(signed))
	err = VerifyPersonalSign(AddrStr, data, rd)
	if nil != err {
		panic("Verify personal error: " + err.Error())
	}

}

func TestParse(t *testing.T) {
	var x = big.NewInt(72727269)
	// var raw = x.Bytes()
	fmt.Println("", hexutil.EncodeBig(x))
}
