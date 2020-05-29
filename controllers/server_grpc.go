package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pangxieke/jwt-go/models"
	"github.com/pangxieke/jwt-go/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type server struct {
	secret []byte
}

func NewServer(secret []byte) *server {
	return &server{secret: secret}
}

func (s *server) checkRequest(req *pb.CreateRequest) error {
	if req.Subject == "" {
		return fmt.Errorf("invalid subject: %s", req.Subject)
	}
	if req.UserId == "" {
		return errors.New("uid is required")
	}
	if req.ExpiredAt != nil && !time.Unix(req.ExpiredAt.Seconds, 0).After(time.Now()) {
		return fmt.Errorf("invalid expired_at, should be later than now")
	}
	return nil
}

func (this *server) Create(ctx context.Context, req *pb.CreateRequest) (token *pb.Token, err error) {
	err = this.checkRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "err = %v", err)
	}
	var payload map[string]interface{}
	if req.Data != "" {
		err = json.Unmarshal([]byte(req.Data), &payload)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid data = %v, err = %v", req.Data, err)
		}
	}

	var exp *time.Time
	if req.ExpiredAt != nil {
		t := time.Unix(req.ExpiredAt.Seconds, 0)
		exp = &t
	}

	jwtToken, err := models.NewJWT(this.secret, req.Subject, req.UserId, exp, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "err = %v", err)
	} else {
		return &pb.Token{Token: jwtToken.Token}, nil
	}
}

func (this *server) Parse(ctx context.Context, req *pb.Token) (*pb.TokenPayload, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token should not be empty")
	}

	token, err := models.ParseJWT(this.secret, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "err = %v", err)
	}

	result := &pb.TokenPayload{}
	dataset := make(map[string]interface{}, 0)
	for k, v := range token.Claims {
		switch k {
		case "subject":
			result.Subject = token.Claims["subject"].(string)
		case "userId":
			result.UserId = token.Claims["userId"].(string)
		case "exp":
			ts := timestamp.Timestamp{Seconds: int64(token.Claims["exp"].(float64))}
			result.ExpiredAt = &ts
		default:
			dataset[k] = v
		}
	}
	b, _ := json.Marshal(dataset)
	result.Data = string(b)
	return result, nil
}

func (this *server) Revoke(ctx context.Context, req *pb.Token) (*empty.Empty, error) {
	return nil, status.Error(codes.Unavailable, "not implemented yet")
}
