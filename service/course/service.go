package course

import (
	"unigo/model/course"
	"unigo/repository"
	"unigo/service"
)

// Service 课程服务
type Service struct {
	*service.CRUDService[course.Course, *repository.CourseRepo]
}

// New 创建课程服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[course.Course, *repository.CourseRepo](repository.NewCourseRepo())}
}
