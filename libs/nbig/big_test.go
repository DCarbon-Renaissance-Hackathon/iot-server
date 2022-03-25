package nbig

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TStruct struct {
	Num *Int `json:"num"`
}

func TestToTwo(t *testing.T) {
	var i, err = NewIntFromString("-0xfffffF")
	if nil != err {
		log.Fatalln("New from string error: ", err)
	}
	fmt.Println("2Byte: ", hexutil.Encode(i.ToTwo(16)))

	var i2 = NewInt(-72727269)
	fmt.Println("32Byte: ", hexutil.Encode(i2.ToTwo(256)))
	raw, err := json.Marshal(&TStruct{Num: i2})
	if nil != err {
		log.Fatalln("Marshal json error: ", err)
	}
	fmt.Println("Raw Json: ", string(raw))

	var st = &TStruct{}
	err = json.Unmarshal(raw, st)
	if nil != err {
		log.Fatalln("Unmarshall json error: ", err)
	}
	log.Println("After unmarshal: ", hexutil.Encode(st.Num.ToTwo(256)))
	log.Println("Test toTwo success")
}
