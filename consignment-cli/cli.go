package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	pb "shippy/consignment-service/proto/consignment"

	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"

	"github.com/micro/go-micro/config/cmd"
)

const (
// ADDRESS           = "localhost:50051"
// DEFAULT_INFO_FILE = "consignment.json"
)

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {
	// service := micro.NewService(micro.Name("go.micro.srv.consignment"))
	// service.Init()
	// client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())

	// 简化微服务创建过程
	cmd.Init()
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// 在命令行中指定新的货物信息 json 文件
	if len(os.Args) < 3 {
		log.Fatalln("Not enough arguments, expecting file and token")
	}
	infoFile := os.Args[1]
	token := os.Args[2]
	log.Printf("infoFile: \n%v\n", infoFile)
	log.Printf("token: \n%v\n\n", token)

	// 解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 创建带有用户token的context,consignment-service取出token，解密用户身份
	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	// 调用 RPC
	// 将货物存储到我们自己的仓库里
	resp, err := client.CreateConsignment(tokenContext, consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	// 新货物是否托运成功
	log.Printf("created: %t", resp.Created)

	// 列出现在所有托运的货物
	resp, err = client.GetConsignments(tokenContext, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments:%v", err)
	}
	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}
}
