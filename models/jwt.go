package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Jwt struct {
	Claims jwt.MapClaims
	Token  string
}

func NewJWT(secret []byte, subject, uid string, expiredAt *time.Time, payload map[string]interface{}) (*Jwt, error) {
	claims := jwt.MapClaims{
		"subject": subject,
		"userId":  uid,
	}
	if expiredAt != nil {
		claims["exp"] = expiredAt.Unix()
	}
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	return &Jwt{Claims: claims, Token: ss}, err
}

func ParseJWT(secret []byte, tokenString string) (*Jwt, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, fmt.Errorf("token parsing failed")
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = fmt.Errorf("That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = fmt.Errorf("Token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				err = fmt.Errorf("Token is not active yet")
			} else {
				err = fmt.Errorf("Could not handle this token, err: %v", err)
			}
			return nil, fmt.Errorf("token invalid, err=%+v", err)
		} else {
			log.Printf("token invalid, err = %+v\n", err)
			return nil, fmt.Errorf("token invalid, unknown error")
		}
	}

	result := &Jwt{Token: tokenString}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result.Claims = claims
	}
	return result, nil
}
