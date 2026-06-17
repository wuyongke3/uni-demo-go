package lecturer

import (
	"unigo/model/lecturer"
	"unigo/repository"
	"unigo/service"
)

// Service 讲师服务 (组合泛型 CRUD 基类)
type Service struct {
	*service.CRUDService[lecturer.Lecturer, *repository.LecturerRepo]
}

// New 创建讲师服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[lecturer.Lecturer, *repository.LecturerRepo](repository.NewLecturerRepo())}
}
