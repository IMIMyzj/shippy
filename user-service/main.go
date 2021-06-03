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

	fmt.Printf("%+v\n", db)
	fmt.Printf("err: %v\n", err)

	defer db.Close()
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	repo := &UserRepository{db}

	// 自动检查 User 结构是否变化
	db.AutoMigrate(&pb.User{})

	// 初始化服务端
	s := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	s.Init()
	t := &TokenService{repo}
	pb.RegisterUserServiceHandler(s.Server(), &handler{repo, t})

	if err := s.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}
