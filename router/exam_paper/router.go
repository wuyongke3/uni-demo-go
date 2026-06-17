package exam_paper

import (
	"unigo/handler"
	"unigo/model/exam_paper"
	svc "unigo/service/exam_paper"

	"github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[exam_paper.ExamPaper]{SVC: svc.New()}

// Register 注册试卷模块路由
func Register(rg *gin.RouterGroup) {
	rg.POST("all", h.All)
	rg.GET("info/:ids", h.Info)
	rg.POST("add", h.Add)
	rg.POST("modify/:id", h.Modify)
	rg.DELETE("delete/:ids", h.Delete)
}
