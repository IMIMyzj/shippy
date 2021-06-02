# Golang微服务

### Golang微服务(一)

#### 使用技术栈

- [x] go

- [x] gRPC

- [x] Makefile

  

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

-----
### Golang微服务(二)

#### 使用技术栈

- [x] go

- [x] go-micro

- [x] Docker

- [x] Makefile

  

#### 整体流程

- <font color=e93b81>**使用.proto定义service 和 message**</font>
  
  - service
    - 定义需要调用的rpc接口
  - message
    - 定义所有rpc接口中输入和输出用到的结构，以及实现它们过程中的结构
- <font color=fc5404> **rpc服务：使用go-micro【需要修改go.mod---->https://www.icode9.com/content-4-729280.html，否则直接按照教程会出错】**</font>
  
  - 不用定义服务端口
  - 仓库接口实现相同
  - 定义微服务的service中把pb.Response返回值放到了函数的入参里
  - main：启动网络监听，启动微服务
    - 微服务注册流程并初始化
    
        ```go 
        server := micro.NewService(…Option)
        server.Init()
        ```
    
    - 为微服务注册数据体
    
      ```go
      pb.RegisterShippingServiceHandler(server.Server(), &service的interface的struct)
      ```
    
    - 启动
    
      ```go
      server.run()
      ```
- <font color=f7a440> **在client的代码中远程调用server内定义的服务**</font>
  
  - 先定义客户端对应的服务并初始化
  
    ```go
    service := micro.NewService(micro.Name("go.micro.srv.consignment"))
    service.Init()
    ```
  
  - 注册为service的客户端
  
    ```go
    client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())
    ```
  
  - 接下去就可以调用相关服务端函数

----

### Golang微服务(三)

#### 使用技术栈

- [x] go
- [x] go-micro
- [x] Docker  | Docker-compose | Dockerfile
- [x] mongoDB  |  postgres



#### 整体流程

- <font color=e93b81>**consignment-service | consignment-cli  | vessel-service**</font>

  - **功能描述**
    - 这三个是一起的，两个服务端一个客户端
    - consignment客户端向consignment服务端请求托运服务，consignment服务端向vessel服务端请求货轮服务
  - **注意点**
    - 依赖的数据库是mongoDB，mongoDB是用database/collection来存储数据的，在docker-compose.yaml中定义了上面的服务以其为依赖，因此会先启动它。打开mongoDB的会话在repository.go中的collection()接口上，直接定位到collection进行增查

- <font color=fc5404> **user-service | user-cli**</font>

  - 功能描述

    - 这两个是一起的，一个服务端一个客户端
    - user-cli向user-service提供个人信息，user-service负责记录信息到postgres数据库中

  - 注意点

    - 依赖的数据库是postgres，默认的用户名/数据库名/密码都是postgres，但是在docker-compose.yaml中自定义了相关的信息

    - **启动postgres特别要注意**：目前不知道原因在哪，需要启动两次postgres，即第一次运行user-service会报错如下，需要再次运行一次【此时postgres已经在运行了】才能正常连接上数据库

      ```shell
      meloneater@meloneater-ubuntu:~/gopath/src/shippy$ docker-compose run user-service
      Creating postgres ... done
      Creating shippy_user-service_run ... done
      Host:database	port:5432	User:userService	Password:12345DbName:userServiceDB
      &{RWMutex:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:0 readerWait:0} Value:<nil> Error:<nil> RowsAffected:0 db:0xc0005ca000 blockGlobalUpdate:false logMode:0 logger:{LogWriter:0xc000189db0} search:<nil> values:{mu:{state:0 sema:0} read:{v:<nil>} dirty:map[] misses:0} parent:0xc0005c00d0 callbacks:0x1e648e0 dialect:0xc0005ac060 singularTable:false nowFuncOverride:<nil>}
      err: dial tcp 172.18.0.5:5432: connect: connection refused
      2021-06-02 12:13:22.045343 I | connect error: dial tcp 172.18.0.5:5432: connect: connection refused
      ERROR: 1
      
      meloneater@meloneater-ubuntu:~/gopath/src/shippy$ docker-compose run user-service
      Creating shippy_user-service_run ... done
      Host:database	port:5432	User:userService	Password:12345DbName:userServiceDB
      &{RWMutex:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:0 readerWait:0} Value:<nil> Error:<nil> RowsAffected:0 db:0xc000600000 blockGlobalUpdate:false logMode:0 logger:{LogWriter:0xc000189db0} search:<nil> values:{mu:{state:0 sema:0} read:{v:<nil>} dirty:map[] misses:0} parent:0xc000606000 callbacks:0x1e648e0 dialect:0xc000604000 singularTable:false nowFuncOverride:<nil>}
      err: <nil>
      2021-06-02 12:13:26.046013 I | Transport [http] Listening on [::]:41921
      2021-06-02 12:13:26.046124 I | Broker [http] Connected to [::]:38709
      2021-06-02 12:13:26.046527 I | Registry [mdns] Registering node: go.micro.srv.user-c0f4d48d-7e77-4438-a42e-1fba51024ca1
      ```

    - **还有一个没解决的问题**：就是micro.Flags的解析，从命令行解析出的数据都是空的，不知道为什么，所以先写死了value
  
- 总体注意点

  - 修改完代码后要make build生成新的二进制文件，再用docker-compose build去利用Dockerfile生成images
  - 在最外层文件夹写了个Makefile是为了方便一些重构和运行操作
