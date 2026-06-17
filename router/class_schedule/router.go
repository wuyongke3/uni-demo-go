package class_schedule

import (
	"unigo/handler"
	"unigo/model/class_schedule"
	svc "unigo/service/class_schedule"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[class_schedule.ClassSchedule]{SVC: svc.New()}

// Register 注册课表模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
