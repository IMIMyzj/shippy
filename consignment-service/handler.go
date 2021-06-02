package main

import (
	"context"
	"log"
	pb "shippy/consignment-service/proto/consignment"
	vesselPb "shippy/vessel-service/proto/vessel"

	"gopkg.in/mgo.v2"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
type handler struct {
	session      *mgo.Session
	vesselClient vesselPb.VesselServiceClient
}

// 从主会话中 Clone() 出新会话处理查询
func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

// 实现接口CreatConsignment
func (h *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	defer h.GetRepo().Close()

	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id
	err = h.GetRepo().Create(req) // 写入数据库
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = req
	return nil
}

// 实现接口CreatConsignment
func (h *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	defer h.GetRepo().Close()
	consignments, err := h.GetRepo().GetAll() // 从数据库查询所有结果
	if err != nil {
		return err
	}
	log.Println("完成所有Consignment数据查询！")
	resp.Consignments = consignments
	return nil
}
