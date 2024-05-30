package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
// 生成包含调用栈信息的字符串-> 了解代码的执行路径
func trace(message string) string {
	var pcs [32]uintptr // 定义调用栈存储数组,存储程序计数器（Program Counter）的值-->调用栈地址
	// 获取调用栈信息
	// n 表示实际获取到的调用栈地址数量
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	// 初始化字符串构建器
	var str strings.Builder
	str.WriteString(message + "\nTRaceback:")
	for _, pc := range pcs[:n] { // 遍历获取到的调用栈地址
		fn := runtime.FuncForPC(pc)                           //  获取对应的函数信息
		file, line := fn.FileLine(pc)                         // 获取调用发生的文件名和行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line)) // 将文件名和行号格式化并添加到字符串构建器中
	}
	return str.String() // 返回构建好的包含调用栈信息的字符串
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil { // 调用 *recover()*，捕获 panic
				message := fmt.Sprintf("%s", err)                               // 将 err 转换为字符串并存储在变量 message
				log.Printf("%s\n\n", trace(message))                            // 将堆栈信息打印在日志中
				c.Fail(http.StatusInternalServerError, "Internal Server Error") // 向用户返回 Internal Server Error
			}
		}()
		c.Next()
	}
}
