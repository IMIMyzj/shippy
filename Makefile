user-cli-rebuild:
	docker container prune
	docker rmi shippy_user-cli
	docker-compose build user-cli

user-service-rebuild:
	docker container prune
	docker rmi shippy_user-service
	docker-compose build user-service

consignment-service-rebuild:
	docker container prune
	docker rmi shippy_consignment-service
	docker-compose build consignment-service

consignment-cli-rebuild:
	docker container prune
	docker rmi shippy_consignment-cli
	docker-compose build consignment-cli

email-service-rebuild:
	docker container prune
	docker rmi shippy_email-service
	docker-compose build email-service

run-user-service:
	docker run --rm --name user-service -net \
	  -e MICRO_ADRESS=":50051" \
	  -e MICRO_REGISTRY="mdns" \
	  -e DB_NAME="userServiceDB" \
	  -e DB_HOST="172.18.0.2" \
	  -e DB_PORT="3306" \
	  -e DB_USER="userService" \
	  -e DB_PASSWORD="12345" \
	  shippy_user-service

run-mysql:
	# 启动mysql
	docker run --rm -d --name mysql -p 3306:3306\
	  -e MYSQL_ROOT_PASSWORD="66666" \
	  -e MYSQL_USER="userService" \
	  -e MYSQL_PASSWORD="12345" \
	  -e MYSQL_DATABASE="userServiceDB" \
	  mysql