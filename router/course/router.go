package course

import (
	"unigo/handler"
	"unigo/model/course"
	svc "unigo/service/course"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[course.Course]{SVC: svc.New()}

// Register 注册课程模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
