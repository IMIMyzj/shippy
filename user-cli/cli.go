package main

import (
	"log"
	"os"

	pb "shippy/user-service/proto/user"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/cmd"
	"golang.org/x/net/context"
)

func main() {
	cmd.Init()

	// 创建 user-service 微服务的客户端
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	// 设置命令行参数:https://learnku.com/docs/go-micro/2.x/write-service/8502
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "name",
				Value: "meloneater",
				Usage: "Your full name",
			},
			cli.StringFlag{
				Name:  "email",
				Value: "meloneater@gmail.com",
				Usage: "Your email",
			},
			cli.StringFlag{
				Name:  "password",
				Value: "12345",
				Usage: "Your password",
			},
			cli.StringFlag{
				Name:  "company",
				Value: "google",
				Usage: "Your company",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")
			company := c.String("company")

			log.Println("得到命令行参数：", name, email, password, company)
			// 创建对应的数据
			r, err := client.Create(context.TODO(), &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
				Company:  company,
			})
			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %+v\n", *(r.User))
			log.Println(*(r.User))

			// 浏览所有数据
			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("Could not list users: %v", err)
			}

			log.Println("所有的用户:")
			for _, v := range getAll.Users {
				log.Printf("用户: %+v\n", *v)
			}

			os.Exit(0)
		}),
	)

	// 启动客户端
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
