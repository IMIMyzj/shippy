package main

import (
	"context"
	"log"
	pb "shippy/consignment-service/proto/consignment"
	vesselPb "shippy/vessel-service/proto/vessel"

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
	repo         Repository
	vesselClient vesselPb.VesselServiceClient
}

// 实现 consignment.pb.go 中的 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
// 将新的货物放到repo.consignments中，然后返回存储情况
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	// 检查是否有合适的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	// 接收承运的货物
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = consignment
	return nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	allConsignments := s.repo.GetAll()
	resp = &pb.Response{Consignments: allConsignments}
	return nil
}

func main() {
	// ----Use go-micro----//
	server := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	server.Init()
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client()) //作为vessel-service的客户端
	repo := Repository{}
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
