package feedback

import (
	"unigo/model/feedback"
	"unigo/repository"
	"unigo/service"
)

// Service 反馈服务
type Service struct {
	*service.CRUDService[feedback.Feedback, *repository.FeedbackRepo]
}

// New 创建反馈服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[feedback.Feedback, *repository.FeedbackRepo](repository.NewFeedbackRepo())}
}
