package class_schedule

import (
	"unigo/model/class_schedule"
	"unigo/repository"
	"unigo/service"
)

// Service 课表服务
type Service struct {
	*service.CRUDService[class_schedule.ClassSchedule, *repository.ClassScheduleRepo]
}

// New 创建课表服务实例
func New() *Service {
	return &Service{CRUDService: service.NewCRUDService[class_schedule.ClassSchedule, *repository.ClassScheduleRepo](repository.NewClassScheduleRepo())}
}
