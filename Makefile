run:
	docker-compose run user-cli command --name="yzj" --email="xxx@gmail.com" --password="12345" --company="google"

user-cli-rebuild:
	docker container prune
	docker rmi shippy_user-cli
	docker-compose build user-cli

user-service-rebuild:
	docker container prune
	docker rmi shippy_user-service
	docker-compose build user-service