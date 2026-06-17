package admin

import (
	"unigo/model/admin"
	"unigo/repository"
	"unigo/service"
)

// Service 管理员服务
type Service struct {
	*service.CRUDService[admin.Admin, *repository.AdminRepo]
}

// New 创建管理员服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[admin.Admin, *repository.AdminRepo](repository.NewAdminRepo())}
}
