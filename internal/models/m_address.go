package models

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Ethereum address.
// Example: 0xdC1A00c3cb7f769ED0C3021A38EC7cfCB5D0631e
type EthAddress string

func (addr *EthAddress) MarshalJSON() ([]byte, error) {
	if nil == addr {
		return nil, nil
	}

	return []byte("\"" + *addr + "\""), nil
}

func (addr *EthAddress) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if string(data) == `""` {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("input for address is invalid")
	}

	var str = string(data[1 : len(data)-1])
	if !common.IsHexAddress(str) {
		return errors.New("input for address is not ethereum address")
	}
	*addr = EthAddress(strings.ToLower(str))

	return nil
}

func (addr *EthAddress) Value() (driver.Value, error) {
	if nil == addr {
		return nil, nil
	}
	return strings.ToLower(string(*addr)), nil
}

func (addr *EthAddress) String() string {
	if nil == addr {
		return ""
	}
	return strings.ToLower(string(*addr))
}

// Data and signed is hex string (0x....)
func (addr *EthAddress) VerifyPersonalSign(data, signed string) error {
	if nil == addr {
		return NewError(ECodeAddressIsEmpty, "Address is empty")
	}

	rawX, err := hexutil.Decode(data)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Data of signature must be hex")
	}

	rawSigned, err := hexutil.Decode(signed)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Signature must be hex")
	}

	err = esign.VerifyPersonalSign(string(*addr), rawX, rawSigned)
	if nil != err {
		return NewError(ECodeInvalidSignature, "Signature invalid")
	}

	return nil
}

func (addr *EthAddress) IsEmpty() bool {
	return nil == addr || *addr == ""
}
