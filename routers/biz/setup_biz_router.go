package biz

import "github.com/gin-gonic/gin"

type Register func(*gin.Engine)

func Init(routers ...Register) *gin.Engine {
	// 注册路由
	rs := append([]Register{}, routers...)

	r := gin.New()
	// 遍历调用方法
	for _, register := range rs {
		register(r)
	}
	return r
}
