package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pangxieke/jwt-go/controllers"
	"github.com/pangxieke/jwt-go/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Println("starting server")
	end := make(chan bool, 1)
	go ServerGRPC()
	go ServerHTTP()
	<-end
}

func loadSecret() ([]byte, error) {
	secret := FetchEnv("SECRET", "")
	if secret == "" {
		return nil, errors.New("Empty secret is not safe! Please set <SECRET>.")
	}
	return hex.DecodeString(secret)
}

func ServerGRPC() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return
	}
	fmt.Println("start grpc:8080")
	secret, err := loadSecret()
	if err != nil {
		log.Fatalf("loading secret failed, err: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTokenServiceServer(s, controllers.NewServer(secret))
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server:%s", err)
	}
}

func FetchEnv(name string, default_value string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return default_value
	}
	return value
}

func ServerHTTP() {
	s := new(controllers.HttpServer)
	http.HandleFunc("/", s.Create)
	http.HandleFunc("/create", s.Create)
	http.HandleFunc("/parse", s.Parse)
	http.HandleFunc("/revoke", s.Revoke)

	fmt.Println("start http:8081")
	http.ListenAndServe("localhost:8081", nil)
}
