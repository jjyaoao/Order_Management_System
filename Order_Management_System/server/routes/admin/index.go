// 命名空间
package admin

// 引入包
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接口函数
func helloWorld(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

// json接口函数
func helloJson(c *gin.Context) {

	// gin框架封装一个json数据
	data := gin.H{
		"name":    "jjyaoao",
		"message": "一个真正的man",
		"age":     19,
	}
	// 传输json数据
	c.JSON(http.StatusOK, data)
}

// 封装json结构体
type create_json struct {
	Name    string `json:"name"`
	Message string
	Age     int
}

// json结构体接口函数
func json_struct(c *gin.Context) {
	// 利用json结构体定义json数据
	data := create_json{
		Name:    "jjyaoao",
		Message: "一个真正的man",
		Age:     19,
	}
	// 传输json数据
	c.JSON(http.StatusOK, data)
}

// 主函数配置路由信息
// 其中*gin.Engine是gin路由的一个实例，在模块的路由中需要引用这个实例
// 注意，要想在不同文件之间调用函数，这个函数首字母就要大写，同时外部调用时也要大写
func Admin(r *gin.Engine) {
	// 定义路由，调用接口函数
	r.GET("/admin", helloWorld)
	// 定义json函数接口
	r.GET("/admin/json", helloJson)
	// 定义json结构体函数接口
	r.GET("/admin/json_struct", json_struct)
}
