package esign

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type TypedDataDomain struct {
	Name              string `json:"name"`              //
	Version           string `json:"version"`           //
	ChainId           int64  `json:"chainid"`           // Hex
	VerifyingContract string `json:"verifyingcontract"` // Address
	Salt              string `json:"salt"`              // Hex
}

type ERC712 struct {
	domain     *TypedDataDomain
	types      *TypedDataField
	domainHash []byte
}

func NewERC712(domain *TypedDataDomain, types *TypedDataField,
) (*ERC712, error) {
	var e712 = &ERC712{
		domain: domain,
		types:  types,
	}
	var data = map[string]interface{}{
		"name":              domain.Name,
		"version":           domain.Version,
		"chainId":           domain.ChainId,
		"verifyingContract": domain.VerifyingContract,
		"salt":              domain.Salt,
	}

	domainHash, err := domainType.Encode(data)
	if nil != err {
		return nil, err
	}
	e712.domainHash = domainHash

	return e712, nil
}

func (e712 *ERC712) Hash(data map[string]interface{},
) ([]byte, error) {
	var dataHash, err = e712.types.Encode(data)
	if nil != err {
		return nil, err
	}
	var sumByte = byteConcat([][]byte{
		{25, 1},
		e712.domainHash,
		dataHash,
	})
	return crypto.Keccak256(sumByte), nil
}

func (e712 *ERC712) Sign(prvStr string, data map[string]interface{}) ([]byte, error) {
	var hash, err = e712.Hash(data)
	if nil != err {
		return nil, err
	}
	return Sign(prvStr, hash)
}
