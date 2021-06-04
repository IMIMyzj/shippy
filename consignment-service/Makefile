build:
# 一定要注意 Makefile 中的缩进，否则执行 make build 可能报错：make: Nothing to be done for `build'.
# protoc 命令前边是一个 Tab，不是四个空格
	protoc --go_out=plugins=micro:. proto/consignment/consignment.proto

# GO编译器生成二进制文件的目标环境
	GOOS=linux GOARCH=amd64 go build

# 根据当前目录下的Dockerfile生成consignment-service镜像
# 将程序执行文件以及所需的run-time环境打包成一个镜像，以后在docker run这个镜像就可以
#	docker build -t consignment-service .

run:
# 在Docker debian容器的50001端口上运行服务
# -p 8080:50051是将容器50051端口映射到本机8080端口
# -e 设置镜像中的环境变量
# mdns作为服务发现的中间层
# 由于MICRO——REGISTRY=mdns和go 1.15及以上的版本会冲突，所以使用1.14.15
	docker run --rm -p 50051:50051 \
	 -e MICRO_SERVER_ADDRESS=:50051 \
	 -e MICRO_REGISTRY=mdns \
	 consignment-service 

