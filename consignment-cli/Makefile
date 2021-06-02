build:
# 生成二进制文件
	GOOS=linux GOARCH=amd64 go build
# 根据Dockerfile生成image
#	docker build -t consignment-cli .

run:
	docker run --rm -e MICRO_REGISTRY=mdns consignment-cli