package esign

import (
	"fmt"
	"testing"
)

func TestPaddingHex(t *testing.T) {
	fmt.Println("", hexPad("0x10000100222", 4))
}

func TestPaddingByte(t *testing.T) {
	fmt.Println(bytePad([]byte{20, 30}, 0))
}

func TestHexConcat(t *testing.T) {
	var arr = []string{"01", "0x02", "0x03", "04"}
	fmt.Println("Hex concat: ", hexConcat(arr))
}

func TestHexPadRight(t *testing.T) {
	var input = "0x112233445566"
	fmt.Println("Hex pad right: ", hexPadRight(input, 5))
}

// 0x11223344556600000000