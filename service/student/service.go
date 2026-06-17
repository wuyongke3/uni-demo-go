package student

import (
	"unigo/model/student"
	"unigo/repository"
	"unigo/service"
)

// Service 学员服务
type Service struct {
	*service.CRUDService[student.Student, *repository.StudentRepo]
}

// New 创建学员服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[student.Student, *repository.StudentRepo](repository.NewStudentRepo())}
}
