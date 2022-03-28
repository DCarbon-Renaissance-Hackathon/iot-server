package esign

import (
	"fmt"
	"testing"
)

func TestPaddingHex(t *testing.T) {
	fmt.Println("", HexPad("0x10000100222", 4))
}

func TestHexPadRight(t *testing.T) {
	var input = "0x112233445566"
	fmt.Println("Hex pad right: ", HexPadRight(input, 5))
}

// 0x11223344556600000000

func TestHexConcat(t *testing.T) {
	var arr = []string{"01", "0x02", "0x03", "04"}
	fmt.Println("Hex concat: ", HexConcat(arr))
}

func TestPaddingByte(t *testing.T) {
	fmt.Println(bytePad([]byte{20, 30}, 0))
}

func TestBytePadRight(t *testing.T) {
	fmt.Println(bytePadRight([]byte{1, 2, 3, 4}, 3))
}

func TestByteConcate(t *testing.T) {
	fmt.Println(byteConcat(
		[][]byte{
			{1, 2, 3},
			{1, 2, 3},
			{4, 5, 6},
		},
	))
}
