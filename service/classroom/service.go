package classroom

import (
	"unigo/model/classroom"
	"unigo/repository"
	"unigo/service"
)

// Service 教室服务
type Service struct {
	*service.CRUDService[classroom.Classroom, *repository.ClassroomRepo]
}

// New 创建教室服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[classroom.Classroom, *repository.ClassroomRepo](repository.NewClassroomRepo())}
}
