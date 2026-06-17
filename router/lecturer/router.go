package lecturer

import (
	"unigo/handler"
	"unigo/model/lecturer"
	svc "unigo/service/lecturer"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[lecturer.Lecturer]{SVC: svc.New()}

// Register 注册讲师模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
