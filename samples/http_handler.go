package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *SampleApp) httpGet(w http.ResponseWriter, req *http.Request) {
	res, err := s.http.Get("http://httpbin.org/get")
	if err != nil {
		responseErr(w, err)
		return
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		responseErr(w, err)
		return
	}
	var jsonBody interface{}
	json.Unmarshal(resBody, &jsonBody)
	responseOk(w, jsonBody)
}
func (s *SampleApp) httpPost(w http.ResponseWriter, req *http.Request) {
	values := map[string]string{"name": "John Doe", "occupation": "gardener"}
	json_data, err := json.Marshal(values)
	if err != nil {
		responseErr(w, err)
		return
	}
	res, err := s.http.Post("http://httpbin.org/post", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		responseErr(w, err)
		return
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		responseErr(w, err)
		return
	}
	var jsonBody interface{}
	json.Unmarshal(resBody, &jsonBody)
	responseOk(w, jsonBody)
}
