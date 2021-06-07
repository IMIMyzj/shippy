package main

import (
	"fmt"
	"log"
	pb "shippy/user-service/proto/user"

	"github.com/micro/go-micro"
)

func main() {
	// 连接到数据库
	db, err := CreateConnection()
	defer db.Close()

	fmt.Printf("%+v\n", db)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}
	repo := &UserRepository{db}

	// 自动检查 User 结构是否变化
	// 如果传入的结构发生变化，那会自动更新结构
	// 结构只增不减
	db.AutoMigrate(&pb.User{})

	// 初始化服务端
	s := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)
	s.Init()

	// 利用topic建立publish端
	publisher := micro.NewPublisher(topic, s.Client())

	t := &TokenService{repo}
	pb.RegisterUserServiceHandler(s.Server(), &handler{repo, t, publisher})

	if err := s.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}
