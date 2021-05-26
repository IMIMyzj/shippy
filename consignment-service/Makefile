build:
# 一定要注意 Makefile 中的缩进，否则执行 make build 可能报错：make: Nothing to be done for `build'.
# protoc 命令前边是一个 Tab，不是四个空格
	protoc --go_out=plugins=grpc:. proto/consignment/consignment.proto