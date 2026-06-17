package feedback

import (
	"unigo/handler"
	"unigo/model/feedback"
	svc "unigo/service/feedback"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[feedback.Feedback]{SVC: svc.New()}

// Register 注册反馈模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
