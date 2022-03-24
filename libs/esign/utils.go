package esign

import "strings"

// Padding to head (left)
func hexPad(input string, numByte int) string {
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

func hexPadRight(input string, numByte int) string {
	if len(input) >= 2 && input[:2] == "0x" {
		input = input[2:]
	}

	if len(input)%2 != 0 {
		input = "0" + input
	}
	var offset = numByte - (len(input)/2)%numByte
	if offset == numByte {
		return "0x" + input
	}
	return "0x" + input + strings.Repeat("00", offset)
}

// Padding to head (left)
func bytePad(input []byte, numByte int) []byte {
	if len(input) < numByte {
		var tmp = make([]byte, numByte-len(input))
		return append(tmp, input...)
	} else if len(input) > numByte {
		return input[:numByte]
	}
	return input
}

func hexConcat(data []string) string {
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
