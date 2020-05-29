package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHttpServer_Create(t *testing.T) {
	assert := assert.New(t)
	reqBody := strings.NewReader(`{
		"subject":"subject",
		"userId": "9527"
	}`)
	resp, err := http.Post("http://localhost:8081/create", "Content-Type:application/json", reqBody)
	assert.Nil(err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	fmt.Println(string(body))

	respRes := struct {
		Token string `json:token`
	}{}
	_ = json.Unmarshal(body, &respRes)
	assert.NotEmpty(respRes.Token)
}

func TestHttpServer_Parse(t *testing.T) {
	assert := assert.New(t)
	reqBody := strings.NewReader(`{
		"subject":"subject",
		"userId": "9527",
		"exp": 1588262400
	}`)
	resp, err := http.Post("http://localhost:8081/create", "Content-Type:application/json", reqBody)
	assert.Nil(err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	//fmt.Println(string(body))

	resp2, err := http.Post("http://localhost:8081/parse", "Content-Type:application/json", bytes.NewReader(body))
	assert.Nil(err)

	body2, err := ioutil.ReadAll(resp2.Body)
	assert.Nil(err)
	fmt.Println(string(body2))

}
