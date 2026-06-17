package exam_paper

import (
	"unigo/model/exam_paper"
	"unigo/repository"
	"unigo/service"
)

// Service 试卷服务
type Service struct {
	*service.CRUDService[exam_paper.ExamPaper, *repository.ExamPaperRepo]
}

// New 创建试卷服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[exam_paper.ExamPaper, *repository.ExamPaperRepo](repository.NewExamPaperRepo())}
}
