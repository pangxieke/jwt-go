package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pangxieke/jwt-go/models"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpServer struct {
	secret []byte
}

func (this *HttpServer) Create(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Subject   string      `json:"subject"`
		UserId    string      `json:"userId"`
		ExpiredAt int64       `json:"exp"`
		Data      interface{} `json"data"`
	}{}

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &req)
	if err != nil {

	}

	//check request
	if req.Subject == "" {
		err = fmt.Errorf("invalid subject: %s", req.Subject)
	}
	if req.UserId == "" {
		err = errors.New("uid is required")
	}
	if err != nil {
		Error(err, w)
	}

	var payload map[string]interface{}

	var exp *time.Time
	if req.ExpiredAt != 0 {
		t := time.Unix(req.ExpiredAt, 0)
		exp = &t
	}

	jwtToken, err := models.NewJWT(this.secret, req.Subject, req.UserId, exp, payload)

	data := make(map[string]interface{})
	data["token"] = jwtToken.Token
	Success(data, w)
}

func (this *HttpServer) Parse(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Token string `json:"token"`
	}{}

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &req)
	if err != nil {

	}
	if req.Token == "" {
		err = fmt.Errorf("token is empty")
		Error(err, w)
	}

	token, err := models.ParseJWT(this.secret, req.Token)
	if err != nil {
		Error(err, w)
	}
	if token == nil {
		err = fmt.Errorf("server error")
		Error(err, w)
	}
	resp := struct {
		Subject   string      `json:"subject"`
		UserId    string      `json:"userId"`
		ExpiredAt int64       `json:"exp"`
		Data      interface{} `json"data"`
	}{}

	if len(token.Claims) != 0 {
		dataset := make(map[string]interface{}, 0)
		for k, v := range token.Claims {
			switch k {
			case "subject":
				resp.Subject = token.Claims["subject"].(string)
			case "userId":
				resp.UserId = token.Claims["userId"].(string)
			case "exp":
				resp.ExpiredAt = int64(token.Claims["exp"].(float64))
			default:
				dataset[k] = v
			}
		}
		resp.Data = dataset
	}

	Success(resp, w)
}

func (this *HttpServer) Revoke(w http.ResponseWriter, r *http.Request) {
	s := "revoke"
	Success(s, w)
}

func Success(data interface{}, w http.ResponseWriter) {
	d, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(d)
}
func Error(err error, w http.ResponseWriter) {

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
