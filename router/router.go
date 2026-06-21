package router

import (
	"net/http"
	"unigo/config"
	"unigo/handler"
	"unigo/middleware"
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

// corsWriter 包装 gin.ResponseWriter，确保每次写入响应时都携带 CORS 头
type corsWriter struct {
	gin.ResponseWriter
	headerWritten bool
}

func (w *corsWriter) Write(data []byte) (int, error) {
	if !w.headerWritten {
		w.setHeader()
		w.headerWritten = true
	}
	return w.ResponseWriter.Write(data)
}

func (w *corsWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.setHeader()
		w.headerWritten = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *corsWriter) setHeader() {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
}

// corsMiddleware CORS 中间件 (包装 ResponseWriter + 处理 OPTIONS 预检)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
			c.Header("Access-Control-Max-Age", "43200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Writer = &corsWriter{ResponseWriter: c.Writer}
		c.Next()
	}
}

// SetupRouter 注册所有路由 (公开登录 + 鉴权保护业务接口)
//
// 路由结构:
//
//	/api/v1
//	├── /auth                    ★ 公开 (无需 Token)
//	│   ├── /lecturer/login     讲师登录
//	│   ├── /lecturer/register  讲师注册
//	│   ├── /student/login      学员登录
//	│   ├── /student/register   学员注册
//	│   └── /me                 获取当前用户信息 (需 Token)
//	│
//	├── /lecturers              ★ 需 Token 鉴权
//	├── /students               ★ 需 Token 鉴权
//	├── /courses                ★ 需 Token 鉴权
//	├── /class-schedules         ★ 需 Token 鉴权
//	├── /classrooms             ★ 需 Token 鉴权
//	├── /exams                  ★ 需 Token 鉴权
//	├── /exam-papers            ★ 需 Token 鉴权
//	├── /feedbacks              ★ 需 Token 鉴权
//	└── /admins                 ★ 需 Token 鉴权
func SetupRouter(r *gin.Engine, cfg config.JWTConfig) {
	api := r.Group("/api/v1")

	// ============================================================
	//  公开路由 (无需 Token)
	// ============================================================
	authGroup := api.Group("/auth")
	authHandler := handler.NewAuthHandler(cfg) // 注入 JWT 配置

	authGroup.POST("/lecturer/login", authHandler.LecturerLogin)       // 讲师登录
	authGroup.POST("/lecturer/register", authHandler.LecturerRegister) // 讲师注册
	authGroup.POST("/student/login", authHandler.StudentLogin)         // 学员登录
	authGroup.POST("/student/register", authHandler.StudentRegister)   // 学员注册

	// ============================================================
	//  需鉴权的业务路由 (所有接口都需要 Bearer Token)
	// ============================================================
	protected := api.Group("") // 使用中间件组保护所有后续路由
	protected.Use(middleware.JWTAuth(cfg))

	// 当前用户信息
	protected.GET("/auth/me", handler.GetCurrentUser)

	// 各模块独立注册路由 (每个子包负责自己的路径和 Handler)
	lecturer.Register(protected.Group("/lecturers"))
	student.Register(protected.Group("/students"))
	course.Register(protected.Group("/courses"))
	class_schedule.Register(protected.Group("/class-schedules"))
	classroom.Register(protected.Group("/classrooms"))
	exam.Register(protected.Group("/exams"))
	exam_paper.Register(protected.Group("/exam-papers"))
	feedback.Register(protected.Group("/feedbacks"))
	admin.Register(protected.Group("/admins"))
}
