package esign

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"regexp"

	"github.com/Dcarbon/libs/nbig"
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

type CBEncode func(value interface{}) ([]byte, error)

type TypedDataField struct {
	// IsArray     bool              `json:"isArray"`
	Name        string            `json:"name"`
	Type        TypedData         `json:"type"`
	Extension   []*TypedDataField `json:"ext"`
	encodeCache CBEncode          `json:"-"`
	domainHash  []byte
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

func (field *TypedDataField) Encode(value interface{}) ([]byte, error) {
	if nil != field.encodeCache {
		err := field.SelectEncodeCb()
		if nil != err {
			return nil, err
		}
	}

	encoded, err := field.encodeCache(value)
	if nil != err {
		log.Printf("Encode field %s error: %s\n", field.Name, err.Error())
		return nil, err
	}
	// fmt.Printf("%s\t\t\t: %s\n", field.Name, hexutil.Encode(encoded))
	return encoded, nil
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

func (field *TypedDataField) encodeAddress(val interface{}) ([]byte, error) {
	switch tVal := val.(type) {
	case string:
		var addrByte, err = hexutil.Decode(tVal)
		if nil != err {
			return nil, err
		}
		return bytePad(addrByte, 32), nil
	case []byte:
		return bytePad(tVal, 32), nil
	default:
		return nil, fmt.Errorf("value for TypedDataField address must be hex string")
	}
}

func (field *TypedDataField) encodeBool(val interface{}) ([]byte, error) {
	var b, ok = val.(bool)
	if !ok {
		return nil, fmt.Errorf("value for TypedDataField bool must be bool")
	}
	if b {
		return bytePad([]byte{1}, 32), nil
	}
	return bytePad([]byte{0}, 32), nil
}

func (field *TypedDataField) encodeBytes(val interface{}) ([]byte, error) {
	var raw, ok = val.([]byte)
	if !ok {
		return nil, fmt.Errorf("value for TypedDataField bytes must be []byte")
	}
	return crypto.Keccak256(raw), nil
}

func (field *TypedDataField) encodeString(val interface{}) ([]byte, error) {
	var raw, ok = val.(string)
	if !ok {
		return nil, fmt.Errorf("value for TypedDataField string must be string")
	}
	return crypto.Keccak256([]byte(raw)), nil
}

func (field *TypedDataField) encodeIntXXX(val interface{}) ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	var ibig = nbig.NewInt(0)
	switch i := val.(type) {
	case int:
		ibig = nbig.NewInt(int64(i))
	case int8:
		ibig = nbig.NewInt(int64(i))
	case int16:
		ibig = nbig.NewInt(int64(i))
	case int32:
		ibig = nbig.NewInt(int64(i))
	case int64:
		ibig = nbig.NewInt(int64(i))
	case uint:
		binary.Write(buf, binary.BigEndian, uint64(i))
		ibig.SetBytes(buf.Bytes())
	case uint8:
		binary.Write(buf, binary.BigEndian, uint64(i))
		ibig.SetBytes(buf.Bytes())
	case uint16:
		binary.Write(buf, binary.BigEndian, uint64(i))
		ibig.SetBytes(buf.Bytes())
	case uint32:
		binary.Write(buf, binary.BigEndian, uint64(i))
		ibig.SetBytes(buf.Bytes())
	case uint64:
		binary.Write(buf, binary.BigEndian, uint64(i))
		ibig.SetBytes(buf.Bytes())
	case string: // Hex
		decoded, err := hexutil.Decode(i)
		if nil != err {
			return nil, err
		}
		buf.Write(decoded)
		ibig.SetBytes(buf.Bytes())
	case *big.Int:
		ibig.Set(i)
	case *nbig.Int:
		ibig.Set(i.Int)
	default:
		return nil, fmt.Errorf("value for TypedDataField Intxx is invalid (%s)", i)
	}

	if ibig == nil {
		return nil, fmt.Errorf("value for TypedDataField Intxx is invalid")
	}
	if ibig.Sign() < 0 && field.Type[0] == 'u' {
		return nil, fmt.Errorf("value is negative for TypedDataField unsign")
	}
	return bytePad(ibig.ToTwo(256), 32), nil
}

func (field *TypedDataField) encodeByteXXX(val interface{}) ([]byte, error) {
	switch i := val.(type) {
	case string:
		return bytePadRight([]byte(i), 32), nil
	case []byte:
		return bytePadRight(i, 32), nil
	}
	return nil, fmt.Errorf("value for TypedDataField string must be string")
}

func (field *TypedDataField) encodeStruct(val interface{}) ([]byte, error) {
	var data, ok = val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value for TypedDataField struct must be map[string]interface")
	}
	var ls = [][]byte{field.domainHash}
	for _, it := range field.Extension {
		var itVal = data[it.Name]
		if nil == itVal {
			return nil, fmt.Errorf("not found value of field %s", it.Name)
		}
		hash, err := it.Encode(itVal)
		if nil != err {
			return nil, err
		}
		ls = append(ls, hash)
	}

	return crypto.Keccak256(byteConcat(ls)), nil
}

func (field *TypedDataField) encodeArray(val interface{}) ([]byte, error) {
	return nil, fmt.Errorf("array encode (TypedDataField) not implement")
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
	field.domainHash = crypto.Keccak256([]byte(domainType))
}
