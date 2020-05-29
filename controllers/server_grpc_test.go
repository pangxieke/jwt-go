package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pangxieke/jwt-go/pb"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"

	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

var testServer *server

func setUp() {
	s := "dc3dc8a96e7053c54ee5267363f9cd803912ea82"
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	testServer = NewServer(b)
}

func TestTokenCreate(t *testing.T) {
	assert := assert.New(t)
	req := pb.CreateRequest{
		Subject: "appkey,app_1234",
		UserId:  "9527",
		ExpiredAt: &timestamp.Timestamp{
			Seconds: time.Now().Add(time.Hour).Unix(),
		},
	}
	res, err := testServer.Create(context.Background(), &req)
	assert.Nil(err)
	assert.NotEqual("", res.Token, res)
	fmt.Printf("res.Token = %+v\n", res.Token)
}

func TestTokenCreate_WithoutSubject(t *testing.T) {
	assert := assert.New(t)
	req := pb.CreateRequest{
		UserId: "9527",
		ExpiredAt: &timestamp.Timestamp{
			Seconds: time.Now().Add(time.Hour).Unix(),
		},
	}
	_, err := testServer.Create(context.Background(), &req)
	assert.NotNil(err)
	assert.Equal(codes.InvalidArgument, status.Code(err))
}

func TestTokenCreate_WithoutUid(t *testing.T) {
	assert := assert.New(t)
	req := pb.CreateRequest{
		Subject: "appkey,app_1234",
		ExpiredAt: &timestamp.Timestamp{
			Seconds: time.Now().Add(time.Hour).Unix(),
		},
	}
	_, err := testServer.Create(context.Background(), &req)
	assert.NotNil(err)
	assert.Equal(codes.InvalidArgument, status.Code(err))
}

func TestTokenCreate_WithInvalidExpiredAt(t *testing.T) {
	assert := assert.New(t)
	req := pb.CreateRequest{
		Subject: "appkey,app_1234",
		UserId:  "9527",
		ExpiredAt: &timestamp.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	_, err := testServer.Create(context.Background(), &req)
	assert.NotNil(err)
	assert.Equal(codes.InvalidArgument, status.Code(err))
}

func TestTokenParse(t *testing.T) {
	assert := assert.New(t)
	req := pb.CreateRequest{
		Subject: "appkey,app_1234",
		UserId:  "9527",
		ExpiredAt: &timestamp.Timestamp{
			Seconds: time.Now().Add(time.Hour).Unix(),
		},
		Data: `
		{
			"i": 123,
			"s": "123"
		}`,
	}
	res, err := testServer.Create(context.Background(), &req)
	assert.Nil(err)
	payload, err := testServer.Parse(context.Background(), &pb.Token{Token: res.Token})
	assert.Nil(err)
	assert.Equal(req.Subject, payload.Subject)
	assert.Equal(req.UserId, payload.UserId)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(payload.Data), &data)
	assert.Nil(err)
	assert.Equal(123.0, data["i"].(float64))
	assert.Equal("123", data["s"].(string))
}
