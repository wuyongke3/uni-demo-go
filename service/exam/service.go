package exam

import (
	"unigo/model/exam"
	"unigo/repository"
	"unigo/service"
)

// Service 考核服务
type Service struct {
	*service.CRUDService[exam.Exam, *repository.ExamRepo]
}

// New 创建考核服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[exam.Exam, *repository.ExamRepo](repository.NewExamRepo())}
}
