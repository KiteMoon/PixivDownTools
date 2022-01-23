package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// 实现一个文件下载函数
func downfile(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("下载文件错误，自动忽略,错误原因:%s", err.Error())
		return
	}
	defer r.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建文件失败，错误如下:%s", err.Error())
		return
	}
	defer file.Close()
	n, err := io.Copy(file, r.Body)
	if err != nil {
		fmt.Printf("写入文件失败，失败原因%s", err.Error())
		return
	}
	fmt.Printf("下载文件成功，文件名%s,写入大小:%d", filename, n)
}
func main() {
	downfile("https://gimg3.baidu.com/search/src=http%3A%2F%2Fpics2.baidu.com%2Ffeed%2F1e30e924b899a901bf8a3b9912c3cc720308f546.jpeg%3Ftoken%3Db9dbba7e142b2ff43b9ea033e489c39c&refer=http%3A%2F%2Fwww.baidu.com&app=2021&size=f360,240&n=0&g=0n&q=75&fmt=auto?sec=1643043600&t=9f44f68e26c90ab46a94921bab9c9922", "ss.png")
}
