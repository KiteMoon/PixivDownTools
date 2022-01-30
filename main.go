package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

// 下载进度条相关，弃用
type DownPorgress struct {
	io.Reader
	NowFileSize    int64
	FinishFileSize int64
}

// 下载进度条,暂时弃用
//func (r *DownPorgress) Read(p []byte) (n int, err error) {
//	n, err = r.Reader.Read(p)
//	r.NowFileSize += int64(n)
//	fmt.Printf("当前下载进度 %f\n", float64(r.NowFileSize*10000/r.FinishFileSize)/100)
//	return
//}

// 实现一个Pixiv文件下载器
func downfile(url, filename string) (code, message, downurl string) {
	// 构建一个http请求downtool，携带pixiv的Referer
	downrequests := http.Client{}
	pixivRequest, _ := http.NewRequest("GET", url, nil)
	//加入header头
	pixivRequest.Header.Add("referer", "https://www.pixiv.net/")
	pixivRequest.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	//发起请求，取得信息
	response, err := downrequests.Do(pixivRequest)
	if err != nil {

		return "500", "下载发生错误，错误原因:" + err.Error(), ""
	}
	if response.StatusCode == 200 {
		file, err := os.Create(filename)
		if err != nil {
			return "500", "创建文件失败，错误如下:" + err.Error(), ""

		}
		defer file.Close()
		n, err := io.Copy(file, response.Body)
		if err != nil {
			return "500", "写入文件失败，错误如下:" + err.Error(), ""

		}
		return "200", "文件下载成功，文件大小" + string(n), "这个功能海米写出来"

	} else if response.StatusCode == 403 { //处理下因为超载导致的问题
		return "403", "读取PIXIV失败，错误如下:" + err.Error(), ""
	}

}

// 实现一个Pixiv链接信息解析
func parsPixivPid(pid string) {

}
func main() {
	server := gin.Default()
	server.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "success",
		})
	})
	server.Run("0.0.0.0:4560")
	//实现一个GIN路由组，该组负责接收下载的文件
	pixivServer := server.Group("/pixiv")
	pixivServer.GET("/get/down/img", func(context *gin.Context) {
		pid := context.Query("pid")

	})
	//for {
	//	var url string
	//	var filename string
	//	fmt.Printf("请输入需要下载的Pixiv图片地址(请提供原始地址):\n")
	//	_, err := fmt.Scan(&url)
	//	if err != nil {
	//		fmt.Printf("发生错误，请输入正确的值")
	//		return
	//	}
	//	fmt.Printf("请输入保存名称:\n")
	//	_, err = fmt.Scan(&filename)
	//	if err != nil {
	//		fmt.Printf("发生错误，请输入正确的值")
	//		return
	//	}
	//	fmt.Println("下载地址:", url)
	//	fmt.Printf("保存文件:%s.png", filename)
	//	downfile(url, (filename + ".png"))
	//}
}
