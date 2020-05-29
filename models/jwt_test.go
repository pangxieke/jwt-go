package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewJWT(t *testing.T) {
	assert := assert.New(t)
	secret := []byte("secret")
	exp := time.Now().Add(time.Hour)
	jwt, err := NewJWT(secret, "appid,appkey", "9527", &exp, nil)
	assert.Nil(err)
	assert.NotEqual("", jwt.Token)

	parsed, err := ParseJWT(secret, jwt.Token)
	assert.Nil(err)
	assert.Equal("appid,appkey", parsed.Claims["subject"].(string))
	assert.Equal("9527", parsed.Claims["userId"].(string))
	assert.Equal(float64(exp.Unix()), parsed.Claims["exp"].(float64))

	exp = time.Now().Add(time.Second)
	jwt, err = NewJWT(secret, "appid,appkey", "uid1", &exp, nil)
	assert.Nil(err)
	assert.NotEqual("", jwt.Token)

	time.Sleep(2 * time.Second)

	parsed, err = ParseJWT(secret, jwt.Token)
	assert.NotNil(err, "token should expire")

}

func TestNewJWT_WithPayloads(t *testing.T) {
	assert := assert.New(t)
	secret := []byte("secret")

	exp := time.Now().Add(time.Hour)
	jwt, err := NewJWT(secret, "appid,appkey", "9527", &exp, map[string]interface{}{
		"integer": 1234,
		"string":  "string",
	})
	assert.Nil(err)
	assert.NotEqual("", jwt.Token)

	parsed, err := ParseJWT(secret, jwt.Token)
	assert.Nil(err)
	assert.Equal(1234.0, parsed.Claims["integer"].(float64))
	assert.Equal("string", parsed.Claims["string"].(string))
}
