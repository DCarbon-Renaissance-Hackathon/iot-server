package models

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Ethereum address
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
