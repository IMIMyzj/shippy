package main

import (
	"context"
	"log"
	pb "shippy/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
	"github.com/pkg/errors"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselReposity struct {
	vessels []*pb.Vessel
}

// 接口实现
func (repo *VesselReposity) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	// 选择最近一条容量和载重都符合的货轮
	for _, v := range repo.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("No vessel can be use")
}

// 定义货船服务
type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, spec *pb.Specification, resp *pb.Response) error {
	v, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

func main() {
	// 写下停留在钢构的货船
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselReposity{vessels}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	server.Init()

	// 注册服务端的API
	pb.RegisterVesselServiceHandler(server.Server(), &service{repo})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
