build:
	protoc --go_out=plugins=micro:. proto/user/user.proto
	GOOS=linux GOARCH=amd64 go build
#	docker build -t user-service .

run:
# 	docker run -p 50053:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns  user-service
# 运行在主机
	docker run --net="host" \
		-p 50052 \
		-e MICRO_SERVER_ADDRESS=:50052 \
		-e MICRO_REGISTRY=mdns \
		-e DISABLE_AUTH=true \
		shippy_consignment-service