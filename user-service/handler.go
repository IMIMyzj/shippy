package main

import (
	"context"
	"errors"
	"log"
	pb "shippy/user-service/proto/user"

	"github.com/micro/go-micro"
	_ "github.com/micro/go-micro/broker/nats"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo         Repository
	tokenService Authable
	Publisher    micro.Publisher
}

const topic = "user.created" // NATS的topic

func (h *handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	// 哈希处理用户输入的密码,定义默认的cost
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPwd)

	// 创建用户, 把哈希密码存到数据库中；用户其他信息都不变
	if err := h.repo.Create(req); err != nil {
		return nil
	}
	resp.User = req

	// 用户注册成功后发布消息，从而建立和email-service的响应
	log.Printf("Called by user-cli to Create user success, now publish event to notify email")
	if err := h.Publisher.Publish(ctx, req); err != nil {
		return err
	}
	return nil
}

func (h *handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	u, err := h.repo.Get(req.Id)
	if err != nil {
		return err
	}
	resp.User = u
	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	users, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = users
	return nil
}

func (h *handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	u, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	// 密码验证，把当前用户的密码和存储着的密码进行比对
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return err
	}
	t, err := h.tokenService.Encode(u) // 此时如果是对的密码，那就对用户的信息进行加密
	if err != nil {
		return err
	}
	resp.Token = t

	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	// 将加密的信息转为claims信息
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}
	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true
	return nil
}
