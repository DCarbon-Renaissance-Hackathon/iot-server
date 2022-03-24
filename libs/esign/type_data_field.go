package esign

import (
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TypedData string

const (
	TypedDataAddress TypedData = "address"
	TypedDataBool    TypedData = "bool"
	TypedDataString  TypedData = "string"
	TypedDataBytes   TypedData = "bytes"
	TypedDataStruct  TypedData = "struct"
)

var regByteXX = regexp.MustCompile(`^byte(\d+)$`)
var regIntXX = regexp.MustCompile(`^(u?)int(\d+)$`)
var regArray = regexp.MustCompile(`^(.*)\[(\d*)\]$`)

var domainType = MustNewTypedDataField(
	"EIP712Domain",
	TypedDataStruct,
	MustNewTypedDataField("name", TypedDataString, nil),
	MustNewTypedDataField("version", TypedDataString, nil),
	MustNewTypedDataField("chainId", "uint256", nil),
	MustNewTypedDataField("verifyingContract", TypedDataAddress, nil),
)

type CBEncode func(value interface{}) (string, error)

type TypedDataField struct {
	// IsArray     bool              `json:"isArray"`
	Name        string            `json:"name"`
	Type        TypedData         `json:"type"`
	Extension   []*TypedDataField `json:"ext"`
	encodeCache CBEncode          `json:"-"`
	domainHash  string
}

func NewTypedDataField(name string, dType TypedData, exts ...*TypedDataField,
) (*TypedDataField, error) {
	var field = &TypedDataField{
		Name:      name,
		Type:      dType,
		Extension: exts,
	}

	var err = field.SelectEncodeCb()
	if dType == TypedDataStruct {
		field.generateDomainHash()
	}

	return field, err
}

func MustNewTypedDataField(name string, dType TypedData, exts ...*TypedDataField,
) *TypedDataField {
	var field, err = NewTypedDataField(name, dType, exts...)
	if nil != err {
		log.Fatalf("Create TypedDataField error %s\n", err.Error())
	}
	return field
}

func (field *TypedDataField) Encode(value interface{}) (string, error) {
	if nil != field.encodeCache {
		err := field.SelectEncodeCb()
		if nil != err {
			return "", err
		}
	}
	return field.encodeCache(value)
}

func (field *TypedDataField) SelectEncodeCb() error {
	if nil != field.encodeCache {
		return nil
	}

	switch field.Type {
	case TypedDataAddress:
		field.encodeCache = field.encodeAddress
		return nil
	case TypedDataBool:
		field.encodeCache = field.encodeBool
		return nil
	case TypedDataBytes:
		field.encodeCache = field.encodeBytes
		return nil
	case TypedDataString:
		field.encodeCache = field.encodeString
		return nil
	case TypedDataStruct:
		field.encodeCache = field.encodeStruct
		return nil
	}

	if regIntXX.Match([]byte(field.Type)) {
		field.encodeCache = field.encodeIntXXX
		return nil
	}

	if regByteXX.Match([]byte(field.Type)) {
		field.encodeCache = field.encodeByteXXX
		return nil
	}

	if regArray.Match([]byte(field.Type)) {
		field.encodeCache = field.encodeArray
		return nil
	}

	return fmt.Errorf("type %s is not support", field.Type)
}

func (field *TypedDataField) encodeAddress(val interface{}) (string, error) {
	var addr, ok = val.(string)
	if !ok {
		return "", fmt.Errorf("value for TypedDataField address must be hex string")
	}
	return hexPad(addr, 32), nil
}

func (field *TypedDataField) encodeBool(val interface{}) (string, error) {
	var b, ok = val.(bool)
	if !ok {
		return "", fmt.Errorf("value for TypedDataField bool must be bool")
	}
	if b {
		return hexPad("0x1", 32), nil
	}
	return hexPad("0x0", 32), nil
}

func (field *TypedDataField) encodeBytes(val interface{}) (string, error) {
	var raw, ok = val.([]byte)
	if !ok {
		return "", fmt.Errorf("value for TypedDataField bytes must be []byte")
	}
	return hexutil.Encode(crypto.Keccak256(raw)), nil
}

func (field *TypedDataField) encodeString(val interface{}) (string, error) {
	var raw, ok = val.(string)
	if !ok {
		return "", fmt.Errorf("value for TypedDataField string must be string")
	}
	return hexutil.Encode(crypto.Keccak256([]byte(raw))), nil
}

func (field *TypedDataField) encodeIntXXX(val interface{}) (string, error) {
	var s = ""
	switch i := val.(type) {
	case int:
		s = strconv.FormatInt(int64(i), 16)
	case int8:
		s = strconv.FormatInt(int64(i), 16)
	case int16:
		s = strconv.FormatInt(int64(i), 16)
	case int32:
		s = strconv.FormatInt(int64(i), 16)
	case int64:
		s = strconv.FormatInt(int64(i), 16)
	case uint:
		s = strconv.FormatInt(int64(i), 16)
	case uint8:
		s = strconv.FormatInt(int64(i), 16)
	case uint16:
		s = strconv.FormatInt(int64(i), 16)
	case uint32:
		s = strconv.FormatInt(int64(i), 16)
	case uint64:
		s = strconv.FormatInt(int64(i), 16)
	case string: // Hex
		s = i
	case big.Int:
		s = hexutil.EncodeBig(&i)
	case *big.Int:
		s = hexutil.EncodeBig(i)
	default:
		return "", fmt.Errorf("value for TypedDataField Intxx is invalid (%s)", i)
	}
	return hexPad(s, 32), nil
}

func (field *TypedDataField) encodeByteXXX(val interface{}) (string, error) {
	switch i := val.(type) {
	case string:
		return hexPadRight(i, 32), nil
	case []byte:
		var hex = hexutil.Encode(i)
		return hexPadRight(hex, 32), nil
	}
	return "", fmt.Errorf("value for TypedDataField string must be string")
}

func (field *TypedDataField) encodeStruct(val interface{}) (string, error) {
	var data, ok = val.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("value for TypedDataField struct must be map[string]interface")
	}
	var ls = []string{field.domainHash}
	for _, it := range field.Extension {
		var itVal = data[it.Name]
		if nil == itVal {
			return "", fmt.Errorf("not found value of field %s", it.Name)
		}
		hash, err := it.Encode(itVal)
		if nil != err {
			return "", err
		}
		ls = append(ls, hash)
	}

	var sum = hexConcat(ls)
	var hashSum = crypto.Keccak256(hexutil.MustDecode(sum))
	return hexutil.Encode(hashSum), nil
}

func (field *TypedDataField) encodeArray(val interface{}) (string, error) {
	return "", fmt.Errorf("not implement")
}

func (field *TypedDataField) generateDomainHash() {
	var domainType = field.Name + "("
	for i, it := range field.Extension {
		domainType += string(it.Type) + " " + it.Name
		if i != len(field.Extension)-1 {
			domainType += ","
		}
	}
	domainType += ")"
	field.domainHash = hexutil.Encode(crypto.Keccak256([]byte(domainType)))
}
