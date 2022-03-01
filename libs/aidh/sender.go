package aidh

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const mimeJSON = "application/json"

//D : dynamic
type D map[string]interface{}

//SendJSON :
func SendJSON(w http.ResponseWriter, code int, data interface{}) {
	var err error
	var raw = []byte("{}")
	if nil != data {
		raw, err = json.Marshal(data)
	}
	if nil != err {
		fmt.Println("Marshall error: ", err)
	} else {
		w.Header().Add("Content-Type", mimeJSON)
		w.WriteHeader(code)
		w.Write(raw)
	}
}

//SendJSONSuccess :
func SendJSONSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, 200, data)
}

//SendJSONErrorBadRequest :
func SendJSONErrorBadRequest(w http.ResponseWriter, msg string) {
	SendJSON(w, 400, map[string]interface{}{
		"msg": msg,
	})
}
