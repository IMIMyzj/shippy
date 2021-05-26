# Golange微服务

### Golang微服务(一)

#### 使用技术栈

- [x] go

- [x] gRPC

  

#### 整体流程

- <font color=e93b81>**使用.proto定义service 和 message**</font>
  - service
    - 定义需要调用的rpc接口
  - message
    - 定义所有rpc接口中输入和输出用到的结构，以及实现它们过程中的结构
- <font color=fc5404> **rpc服务：在server的代码中实现.proto的service并且启动grpc服务**</font>
  - 定义服务端口
  - 定义仓库接口并实现
  - 定义微服务service，所有sevice的第一个参数为 ctx context.Context，用于保存上下文信息
  - main：启动网络监听，启动微服务
    - 启动网络监听：`listener, err := net.Listen("tcp", PORT)`
    - 启动grpc服务端：`server := grpc.NewServer()`
    -  为微服务注册数据体：`pb.RegisterShippingServiceServer(server, &service的interface的struct)`
    -  开启监听：`server.Serve(listener)`
- <font color=f7a440> **在client的代码中远程调用server内定义的服务**</font>
  - 连接到服务端口：`conn, err := grpc.Dial(ADDRESS,  grpc.WithInsecure)`
  - 初始化grpc客户端：`client := pb.NewShippingServiceClient(conn)`
  - 远程调用服务：`client.service_func(xxxx)`
