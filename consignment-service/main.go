package main

import (
	"context"
	"errors"
	"log"
	"os"
	pb "shippy/consignment-service/proto/consignment"
	userPb "shippy/user-service/proto/user"
	vesselPb "shippy/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

const (
	DEFAULT_HOST = "localhost:27017"
)

func main() {
	// 连接MongoDB数据库
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = DEFAULT_HOST
	}

	session, err := CreateSession(dbHost)
	defer session.Close()
	if err != nil {
		log.Fatalf("Create session error: %v\n", err)
	}

	// 启动微服务
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)
	server.Init()
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client()) //作为vessel-service的客户端
	pb.RegisterShippingServiceHandler(server.Server(), &handler{session, vClient})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// AuthWrapper是高阶函数，入参：下一步函数，出参：认证函数
// 在返回的内部函数处理完后，进入入参函数处理
// token从consignment-cli上下文中取出，再到user-service中做认证
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		// consignment-service独立测试时不进行认证
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx) // 获取cli.go的context中存入的token
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		token := meta["Token"]
		log.Printf("token: %v\n", token) // 显示下密钥，可以不加这句

		// Auth
		authClient := userPb.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(context.Background(), &userPb.Token{
			Token: token,
		})
		log.Println("Auth Resp:", authResp)
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
