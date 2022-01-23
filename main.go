package main

import (
	"fmt"
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

// 实现一个文件下载函数
func downfile(url, filename string) {
	// 构建一个http请求downtool，携带pixiv的Referer
	downrequests := http.Client{}
	pixivRequest, _ := http.NewRequest("GET", url, nil)
	//加入header头
	pixivRequest.Header.Add("referer", "https://www.pixiv.net/")
	pixivRequest.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	//发起请求，取得信息
	response, err := downrequests.Do(pixivRequest)
	if err != nil {
		fmt.Println("下载发生错误，错误原因:", err)
		return
	}
	if response.StatusCode == 200 {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("创建文件失败，错误如下:%s\n", err.Error())
			return
		}
		defer file.Close()
		n, err := io.Copy(file, response.Body)

		if err != nil {
			fmt.Printf("写入文件失败，失败原因%s\n", err.Error())
			return
		}
		fmt.Printf("下载文件成功，文件名%s,写入大小:%d\n", filename, n)
	} else if response.StatusCode == 403 { //处理下因为超载导致的问题
		fmt.Printf("下载文件失败，Pixiv返回403\n")
	}

}

func main() {
	for {
		var url string
		var filename string
		fmt.Printf("请输入需要下载的Pixiv图片地址(请提供原始地址):\n")
		_, err := fmt.Scan(&url)
		if err != nil {
			fmt.Printf("发生错误，请输入正确的值")
			return
		}
		fmt.Printf("请输入保存名称:\n")
		_, err = fmt.Scan(&filename)
		if err != nil {
			fmt.Printf("发生错误，请输入正确的值")
			return
		}
		fmt.Println("下载地址:", url)
		fmt.Printf("保存文件:%s.png", filename)
		downfile(url, (filename + ".png"))
	}
}
