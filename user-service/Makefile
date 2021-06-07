build:
	protoc --go_out=plugins=micro:. proto/user/user.proto
	GOOS=linux GOARCH=amd64 go build
#	docker build -t user-service .

run_db:
	docker run --rm -d --name postgres -p 5432:5432 \
		-e POSTGRES_DB="userServiceDB" \
		-e POSTGRES_USER="userService" \
		-e POSTGRES_PASSWORD="12345" \
		postgres

run:
# 	docker run -p 50053:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns  user-service
# 运行在主机
	docker run --rm \
		--link email-service:db \
		-e MICRO_REGISTRY=mdns \
		-e DB_NAME="userServiceDB" \
        -e DB_HOST="database" \
        -e DB_PORT="5432" \
        -e DB_USER="userService" \
        -e DB_PASSWORD="12345" \
		-e MICRO_BROKER=nats \
		-e MICRO_BROKER_ADDRESS=0.0.0.0:4222 \
		shippy_user-service