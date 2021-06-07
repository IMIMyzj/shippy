package main

import (
	"log"
	pb "shippy/user-service/proto/user"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	Get(id string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
	Create(*pb.User) error
	GetByEmail(email string) (*pb.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

// 根据用户ID找到信息
func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var u *pb.User
	u.Id = id
	if err := repo.db.First(&u).Error; err != nil {
		return nil, err
	}
	log.Printf("MySQL找到用户:%+v\n!", *u)
	return u, nil
}

// 根据一堆用户ID找到信息
func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	log.Printf("MySQL找到用户如下：\n")
	for _, v := range users {
		log.Printf("用户: %+v\n", *v)
	}
	return users, nil
}

// 创建数据库
func (repo *UserRepository) Create(u *pb.User) error {
	if err := repo.db.Create(&u).Error; err != nil {
		return err
	}
	log.Printf("MySQL数据库用户%+v创建成功！", *u)
	return nil
}

// 根据email查找用户,不能传入密码数据，否则会被改成数据库内部的，对比哈希密码的时候就会对比相同的string而报错
func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &pb.User{}

	if err := repo.db.Where("email = ?", email).Find(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
