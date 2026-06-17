package classroom

import (
	"unigo/handler"
	"unigo/model/classroom"
	svc "unigo/service/classroom"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[classroom.Classroom]{SVC: svc.New()}

// Register 注册教室模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
