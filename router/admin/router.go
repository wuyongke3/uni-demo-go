package admin

import (
	"unigo/handler"
	"unigo/model/admin"
	svc "unigo/service/admin"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[admin.Admin]{SVC: svc.New()}

// Register 注册管理员模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
