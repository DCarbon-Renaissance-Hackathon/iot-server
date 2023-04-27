package mids

import (
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

var key = utils.StringEnv("JWT_KEY", "")

func TestDecode(t *testing.T) {
	var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMxNDQ0MzMsImlkIjoxLCJyb2xlIjoiYWRtaW4iLCJldGgiOiIweDU3ZDdENzJmNTRCOERiZDg2NjA2MGIwQ0YyNjVhMmE0NWU4ZWY3MmIifQ.yOsoUZ_FZun9gsEFyBvBgjAYaHkq4H87iH9ZrFZhhR8"

	data, err := DecodeJWT(key, token)
	utils.PanicError("", err)
	utils.Dump("", data)

}

func TestA2M(t *testing.T) {
	token, err := EncodeJWT(key, &models.User{
		ID:   1,
		Role: "super-admin",
	}, 10*365*86400)
	utils.PanicError("", err)
	utils.Dump("Token ", token)
}
