package main

import (
	"fmt"
	"log"
	"os"
	pb "shippy/consignment-service/proto/consignment"
	vesselPb "shippy/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
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
	fmt.Printf("MongoDB host: %s\n", dbHost)
	session, err := CreateSession(dbHost)
	defer session.Close()
	if err != nil {
		log.Fatalf("Create session error: %v\n", err)
	}

	// 启动微服务
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	server.Init()
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client()) //作为vessel-service的客户端
	pb.RegisterShippingServiceHandler(server.Server(), &handler{session, vClient})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
