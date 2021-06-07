package main

import (
	"context"
	"log"
	userPb "shippy/user-service/proto/user"

	"github.com/micro/go-micro"
	_ "github.com/micro/go-micro/broker/nats"
)

const topic = "user.created" // 设置topic

type Subscriber struct{}

func main() {
	// 注册邮件服务
	srv := micro.NewService(
		micro.Name("go.micro.srv.email"),
		micro.Version("latest"),
	)
	srv.Init()

	// 注册发布订阅服务
	micro.RegisterSubscriber(topic, srv.Server(), new(Subscriber))

	if err := srv.Run(); err != nil {
		log.Fatalf("srv run error: %v\n", err)
	}
}

// 获取到从user-cli的handle中Create发送过来的user信息
func (sub *Subscriber) Process(ctx context.Context, user *userPb.User) error {
	log.Println("[Picked up a new message]")
	log.Println("[Sending email to]:", user.Name)
	return nil
}
