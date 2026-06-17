package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})

	v1 := r.Group("/v1")
	{
		// 通过 localhost:8080/v1/hello访问，以此类推
		v1.GET("/hello", func(c *gin.Context) {
			v1.GET("/world", func(c *gin.Context) {
				c.String(http.StatusOK, "hello World!")
			})
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello")
		v2.GET("/world")
	}

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
