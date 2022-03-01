package utils

import (
	"encoding/json"
	"log"
)

//PanicError :
func PanicError(msg string, err error) {
	if nil != err {
		panic(msg + " error: " + err.Error())
	}
}

//Dump : use for test
func Dump(label string, data interface{}) {
	if nil == data {
		log.Println(label+" dump: ", nil)
		return
	}
	var raw, err = json.MarshalIndent(data, "", "  ")
	if nil != err {
		panic(err)
	}
	log.Println(label + " dump: " + string(raw))
}
