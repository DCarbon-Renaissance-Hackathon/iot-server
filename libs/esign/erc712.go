package esign

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
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

func MustNewERC712(domain *TypedDataDomain, types *TypedDataField,
) *ERC712 {
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
		panic(err)
	}

	e712.domainHash = domainHash
	return e712
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
	// fmt.Println("Prefix: ", hexutil.Encode([]byte{25, 1}))
	// fmt.Println("Domain hash: ", hexutil.Encode(e712.domainHash))
	// fmt.Println("Hashed struct: ", hexutil.Encode(dataHash))
	return crypto.Keccak256(sumByte), nil
}

func (e712 *ERC712) Sign(prvStr string, data map[string]interface{}) ([]byte, error) {
	prv, err := crypto.HexToECDSA(prvStr)
	if nil != err {
		return nil, err
	}

	hash, err := e712.Hash(data)
	if nil != err {
		return nil, err
	}
	// fmt.Println("Hashed struct: ", hexutil.Encode(hash))

	signed, err := crypto.Sign(hash[:], prv)
	if nil != err {
		return nil, err
	}
	signed[64] += 27
	return signed, nil

}

func (e712 *ERC712) Verify(addrStr string, signed []byte, data map[string]interface{}) error {
	inputPub := common.HexToAddress(addrStr)

	hash, err := e712.Hash(data)
	if nil != err {
		return err
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

	if crypto.PubkeyToAddress(*pubkey) == inputPub {
		return nil
	}
	return errors.New("not match")
}
