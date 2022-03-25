package nbig

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var regString = regexp.MustCompile(`"*"$`)
var regHexJson = regexp.MustCompile(`^(-?)0x[0-9a-zA-Z]*$`)

type Int struct {
	*big.Int
}

func NewInt(val int64) *Int {
	var i = &Int{
		Int: big.NewInt(val),
	}
	return i
}

//  string must be hex
func NewIntFromString(s string) (*Int, error) {
	return newIntFromString(s)
}

func (i *Int) ToTwo(width int) []byte {
	var nByte = width / 8
	if i.Sign() >= 0 {
		var current = i.Bytes()
		if width >= nByte {
			return current[:nByte]
		}
		return append(make([]byte, nByte-len(current)), current...)
	}

	var abs = big.NewInt(0)
	abs.SetBytes(i.Bytes())

	var max = big.NewInt(0)
	max.SetBytes(bytes.Repeat([]byte{0xff}, nByte))

	max.Sub(max, abs).Add(max, big.NewInt(1))

	return max.Bytes()
}

func (i *Int) MarshalJSON() ([]byte, error) {
	var hex = hexutil.Encode(i.Bytes())
	if i.Sign() < 0 {
		hex = "-" + hex
	}
	hex = fmt.Sprintf(`"%s"`, hex)
	return []byte(hex), nil
}

func (i *Int) UnmarshalJSON(data []byte) error {
	var s = string(data)
	if nil == i {
		i = new(Int)
	}

	if i.Int == nil {
		i.Int = big.NewInt(0)
	}

	if !regString.Match(data) {
		i64, err := strconv.ParseInt(s, 10, 64)
		if nil != err {
			return err
		}
		i.Int = big.NewInt(i64)
		return nil
	}
	s = s[1 : len(s)-1]

	iTmp, err := newIntFromString(s)
	if nil != err {
		return errors.New("input for big number must be hex string or int64")
	}
	i.Int = iTmp.Int

	return nil
}

func newIntFromString(input string) (*Int, error) {
	var i = &Int{
		Int: big.NewInt(0),
	}

	var isNeg = false
	if regHexJson.Match([]byte(input)) {
		if input[0] == '-' {
			isNeg = true
			input = input[1:]
		}
		abs, err := hexutil.Decode(input)
		if nil != err {
			return nil, err
		}

		i.SetBytes(abs)
		if isNeg {
			i.Mul(i.Int, big.NewInt(-1))
		}
		return i, nil
	}
	return nil, errors.New("input for big number must be hex string")
}
