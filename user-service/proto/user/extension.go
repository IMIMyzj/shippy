package user

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	// uuid, err := uuid.NewV4()
	// if err != nil {
	// 	log.Fatalf("created uuid error: %v\n", err)
	// }

	uuid := uuid.NewV4()
	return scope.SetColumn("Id", uuid.String())
}
