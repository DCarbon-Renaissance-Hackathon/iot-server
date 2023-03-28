package mids

import (
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
)

func TestDecode(t *testing.T) {
	var key = "047120jlcvndlj092u3jrlnvldnbvlajw021981yu0rhvklndvlkhsr921y49"
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMxNDQ0MzMsImlkIjoxLCJyb2xlIjoiYWRtaW4iLCJldGgiOiIweDU3ZDdENzJmNTRCOERiZDg2NjA2MGIwQ0YyNjVhMmE0NWU4ZWY3MmIifQ.yOsoUZ_FZun9gsEFyBvBgjAYaHkq4H87iH9ZrFZhhR8"

	data, err := DecodeJWT(key, token)
	utils.PanicError("", err)
	utils.Dump("", data)

}

func TestA2M(t *testing.T) {

}
