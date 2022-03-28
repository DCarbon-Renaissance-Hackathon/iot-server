package esign

import "strings"

// Padding to head (left)
func HexPad(input string, numByte int) string {
	if len(input) >= 2 && input[:2] == "0x" {
		input = input[2:]
	}
	var sz = len(input)
	if sz < numByte*2 {
		input = strings.Repeat("0", numByte*2-sz) + input
	} else if sz > numByte*2 {
		input = input[:2*numByte]
	}
	return "0x" + input
}

func HexPadRight(input string, numByte int) string {
	if len(input) >= 2 && input[:2] == "0x" {
		input = input[2:]
	}

	if len(input)%2 != 0 {
		input = "0" + input
	}
	var offset = numByte - (len(input)/2)%numByte
	if offset == numByte && len(input) > 0 {
		return "0x" + input
	}
	return "0x" + input + strings.Repeat("00", offset)
}

func HexConcat(data []string) string {
	var rs = "0x"
	for _, it := range data {
		if len(it) >= 2 && it[:2] == "0x" {
			rs += it[2:]
		} else {
			rs += it
		}
	}
	return rs
}

// Padding to head (left)
func bytePad(input []byte, numByte int) []byte {
	if len(input) < numByte {
		return append(make([]byte, numByte-len(input)), input...)
	}

	if len(input) > numByte {
		return input[:numByte]
	}
	return input
}

func bytePadRight(input []byte, numByte int) []byte {
	var offset = numByte - (len(input) % numByte)
	if offset == numByte && len(input) > 0 {
		return input
	}
	return append(input, make([]byte, offset)...)
}

func byteConcat(input [][]byte) []byte {
	var rs = make([]byte, 0)
	for _, it := range input {
		rs = append(rs, it...)
	}
	return rs
}
