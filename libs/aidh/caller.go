package aidh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

//CallGet : call http with method get (type: json)
func CallGet(URL string, data interface{}) error {
	var resp, err = http.Get(URL)
	if nil != err {
		return err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return err
	}

	if resp.StatusCode != 200 {
		var e = &Error{}
		json.Unmarshal(raw, e)
		return e
	}

	if nil != data {
		err = json.Unmarshal(raw, data)
	}
	return err
}

//Caller :
type Caller struct {
	client       *http.Client
	headersMutex sync.RWMutex
	headers      map[string]string
	errMaker     func() error
}

//NewCaller :
func NewCaller(headers map[string]string) *Caller {
	var c = &Caller{
		headers: headers,
		client:  &http.Client{},
	}
	return c
}

func (c *Caller) SetErrorMaker(maker func() error) {
	c.errMaker = maker
}

//Get : http get
func (c *Caller) Get(URL string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if nil != err {
		return err
	}
	// fmt.Println("GET: ", URL)
	// c.printHeader()
	c.addRequestHeader(req)

	resp, err := c.client.Do(req)
	if nil != err {
		return err
	}
	return c.parseResponse(resp, data)
}

//Post :
func (c *Caller) Post(URL string, payload, data interface{}) error {
	var buff *bytes.Buffer
	if payload != nil {
		raw, err := json.Marshal(payload)
		if nil != err {
			return err
		}
		// log.Println("Payload: ", string(raw))
		buff = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(http.MethodPost, URL, buff)
	if nil != err {
		return err
	}
	c.addRequestHeader(req)

	resp, err := c.client.Do(req)
	if nil != err {
		return err
	}
	return c.parseResponse(resp, data)
}

//Post :
func (c *Caller) PostWithResponse(URL string, payload, data interface{}) (*http.Response, error) {
	var buff *bytes.Buffer
	if payload != nil {
		raw, err := json.Marshal(payload)
		if nil != err {
			return nil, err
		}
		buff = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(http.MethodPost, URL, buff)
	if nil != err {
		return nil, err
	}
	c.addRequestHeader(req)

	resp, err := c.client.Do(req)
	if nil != err {
		return nil, err
	}
	return resp, c.parseResponse(resp, data)
}

//Put : put method
func (c *Caller) Put(URL string, payload, data interface{}) error {
	var buff *bytes.Buffer
	if payload != nil {
		raw, err := json.Marshal(payload)
		if nil != err {
			return err
		}
		buff = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(http.MethodPut, URL, buff)
	if nil != err {
		return err
	}
	c.addRequestHeader(req)

	resp, err := c.client.Do(req)
	if nil != err {
		return err
	}
	return c.parseResponse(resp, data)
}

func (c *Caller) FormFile(URL string, payload FormFields, data interface{}) error {
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)

	for _, field := range payload {
		if field.Type == FormFieldText {
			writer.WriteField(field.Key, field.Value)
		} else if field.Type == FormFieldFile {
			f, err := os.Open(field.Value)
			if nil != err {
				return err
			}
			defer f.Close()

			part, err := writer.CreateFormFile(field.Key, filepath.Base(f.Name()))
			if nil != err {
				return err
			}
			_, err = io.Copy(part, f)
			if nil != err {
				return err
			}

		}
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, URL, buff)
	if nil != err {
		return err
	}
	c.addRequestHeader(req)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := c.client.Do(req)
	if nil != err {
		return err
	}
	return c.parseResponse(resp, data)
}

//AddHeader :
func (c *Caller) AddHeader(k, v string) {
	c.headersMutex.Lock()
	defer c.headersMutex.Unlock()

	c.headers[k] = v
}

//RemoveHeader :
func (c *Caller) RemoveHeader(k string) {
	c.headersMutex.Lock()
	defer c.headersMutex.Unlock()
	delete(c.headers, k)
}

func (c *Caller) addRequestHeader(req *http.Request) {
	c.headersMutex.RLock()
	defer c.headersMutex.RUnlock()

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
}

func (c *Caller) parseResponse(resp *http.Response, data interface{}) error {
	if nil == data {
		return nil
	}
	defer func() {
		if nil != resp.Body {
			resp.Body.Close()
		}
	}()
	raw, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return err
	}

	// if len(raw) < 200 {
	// 	fmt.Println("Raw http: ", string(raw))
	// } else {
	// 	fmt.Println("Raw http too big")
	// }

	if resp.StatusCode >= 300 {
		if c.errMaker != nil {
			var e = c.errMaker()
			json.Unmarshal(raw, e)
			return e
		}
		// return fmt.Errorf(
		// 	"status code:%d url:%s body:%s",
		// 	resp.StatusCode,
		// 	resp.Request.RequestURI,
		// 	string(raw),
		// )
		return fmt.Errorf(
			string(raw),
		)
	}

	return json.Unmarshal(raw, data)
}

// func (c *Caller) printHeader() {
// 	c.headersMutex.RLock()
// 	defer c.headersMutex.RUnlock()

// 	for k, v := range c.headers {
// 		fmt.Println("\t" + k + ": " + v)
// 	}
// }
