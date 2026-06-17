package exam

import (
	"unigo/handler"
	"unigo/model/exam"
	svc "unigo/service/exam"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[exam.Exam]{SVC: svc.New()}

// Register 注册考核模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
