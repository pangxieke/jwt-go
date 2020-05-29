package main

import (
	"context"
	"fmt"
	"github.com/pangxieke/jwt-go/pb"
	"google.golang.org/grpc"
	"time"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("err")
	}

	client := pb.NewTokenServiceClient(conn)
	req := pb.CreateRequest{
		Subject: "abc",
		UserId:  "9527",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Create(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(resp)
	token := resp.Token
	fmt.Println("token:", token)

	req2 := pb.Token{
		Token: token,
	}
	ctx, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	resp2, err := client.Parse(ctx, &req2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp2)

}
