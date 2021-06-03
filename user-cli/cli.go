package main

import (
	"log"
	"os"

	pb "shippy/user-service/proto/user"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/cmd"
	"golang.org/x/net/context"
)

func main() {
	cmd.Init()

	// 创建 user-service 微服务的客户端
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	// 写死用户信息
	name := "meloneater"
	email := "meloneater@gmail.com"
	password := "12345"
	company := "zju"

	// 创建对应的数据
	resp, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %+v\n", *(resp.User))

	// 浏览所有数据
	allResp, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	log.Println("所有的用户:")
	for _, v := range allResp.Users {
		log.Printf("用户: %+v\n", *v)
	}

	// 根据邮箱和密码进行用户查寻
	authResp, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("auth failed: %+v", err)
	}
	log.Println("token: ", authResp.Token)

	// 退出
	os.Exit(0)
}
