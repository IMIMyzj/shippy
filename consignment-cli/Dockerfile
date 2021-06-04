FROM debian:latest
RUN mkdir -p /app
WORKDIR /app

# 把当前的货物信息文件都拷贝到/app中
ADD consignment.json /app/consignment.json
ADD consignment-cli /app/consignment-cli

ENTRYPOINT ["./consignment-cli"]
CMD [ "./consignment-cli" ]