# 使用最新的debian作为基础镜像
FROM debian:latest

# 在容器目录下茶创建app目录
RUN mkdir /app

# 切换工作目录到/app下
WORKDIR /app

# 将服务端运行文件copy到/app下
ADD consignment-service /app/consignment-service

# 运行服务端
CMD ["./consignment-service"]