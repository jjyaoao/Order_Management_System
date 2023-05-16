package routes

import "github.com/gin-gonic/gin"

// 定义接收路由的形式（类型），为gin引擎模板
type Option func(*gin.Engine)

// 定义一个存放多个路由的数组，数组内形式为gin框架引擎模板
// 即上方定义的形式
var options = []Option{}

// 定义一个方法，用于将其他文件引入的路由接口放进数组
// 其中...是go语言中的一种语法糖，用于接收不确定数量的参数
// 函数中的两个值为 变量名 变量类型
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 定义初始化合成引擎模板的函数，用于将多个
// 引擎模板合成为一个新的引擎模板。
func Arrange() *gin.Engine {
	r := gin.New()
	for _, opt := range options {
		opt(r)
	}
	return r
}
