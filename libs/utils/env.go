package utils

import (
	"log"
	"os"
	"strconv"
)

//StringOrEnv :
func StringEnv(key string, dval string) string {
	if os.Getenv(key) == "" {
		return dval
	}
	return os.Getenv(key)
}

//IntOrEnv :
func IntEnv(key string, dval int) int {
	if os.Getenv(key) == "" {
		return dval
	}
	var ev, err = strconv.Atoi(os.Getenv(key))
	if nil == err {
		log.Printf("Env %s must be number. Current val:%s\n", key, os.Getenv(key))
		return dval
	}
	return ev
}

//Int64OrEnv :
func Int64Env(key string, dval int64) int64 {
	if os.Getenv(key) == "" {
		return dval
	}
	var ev, err = strconv.ParseInt(os.Getenv(key), 10, 64)
	if nil == err {
		log.Printf("Env %s must be number. Current val:%s\n", key, os.Getenv(key))
		return dval
	}
	return ev
}

//BoolOrEnv :
func BoolEnv(key string, dval bool) bool {
	if os.Getenv(key) == "" {
		return dval
	}
	if os.Getenv(key) == "none" {
		return false
	}
	return true
}
