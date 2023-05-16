// 命名空间
package web

// 引入包
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接口函数
func helloWorld(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

// 主函数配置路由信息
// 其中*gin.Engine是gin路由的一个实例，在模块的路由中需要引用这个实例
// 注意，要想在不同文件之间调用函数，这个函数首字母就要大写，同时外部调用时也要大写
func Web(r *gin.Engine) {
	r.GET("/web", helloWorld)
}
