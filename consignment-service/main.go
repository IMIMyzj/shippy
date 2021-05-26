package main

import (
	"context"
	"log"
	pb "shippy/consignment-service/proto/consignment"

	"github.com/micro/go-micro"
)

const (
	PORT = ":50051"
)

// ------仓库接口------ //
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
	GetAll() []*pb.Consignment                                   // 获取仓库中所有的货物
}

// 我们存放多批货物的仓库，实现了 IRepository 接口
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// -----定义微服务-----//
type service struct {
	repo Repository
}

// 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
// 将新的货物放到repo.consignments中，然后返回存储情况
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp = &pb.Response{Created: true, Consignment: consignment}
	return nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	resp = &pb.Response{Consignments: allConsignments}
	return nil
}

func main() {
	// ----Use grpc----//
	// server :=
	// listener, err := net.Listen("tcp", PORT)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// log.Printf("listen on: %s\n", PORT)

	// server := grpc.NewServer()
	//--------------------//

	// ----Use go-micro----//
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	server.Init()
	//--------------------//

	repo := Repository{}
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo})

	// ----Use grpc----//
	// pb.RegisterShippingServiceServer(server, &service{repo})
	// if err := server.Serve(listener); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	// ---------------//

	// ----Use go-micro----//
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//--------------------//
}
