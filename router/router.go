package router

import (
	"unigo/router/admin"
	"unigo/router/class_schedule"
	"unigo/router/classroom"
	"unigo/router/course"
	"unigo/router/exam"
	"unigo/router/exam_paper"
	"unigo/router/feedback"
	"unigo/router/lecturer"
	"unigo/router/student"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册所有业务路由
//
// 路由结构:
//
//	/api/v1
//	├── /lecturers        讲师模块
//	├── /students         学员模块
//	├── /courses          课程模块
//	├── /class-schedules  课表模块
//	├── /classrooms       教室模块
//	├── /exams            考核模块
//	├── /exam-papers      试卷模块
//	├── /feedbacks        反馈模块
//	└── /admins           管理员模块
func SetupRouter(r *gin.Engine) {
	api := r.Group("/api/v1")

	// 各模块独立注册路由 (每个子包负责自己的路径和 Handler)
	lecturer.Register(api.Group("/lecturers"))
	student.Register(api.Group("/students"))
	course.Register(api.Group("/courses"))
	class_schedule.Register(api.Group("/class-schedules"))
	classroom.Register(api.Group("/classrooms"))
	exam.Register(api.Group("/exams"))
	exam_paper.Register(api.Group("/exam-papers"))
	feedback.Register(api.Group("/feedbacks"))
	admin.Register(api.Group("/admins"))
}
