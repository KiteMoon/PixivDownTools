package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// 下载进度条,暂时弃用
//func (r *DownPorgress) Read(p []byte) (n int, err error) {
//	n, err = r.Reader.Read(p)
//	r.NowFileSize += int64(n)
//	fmt.Printf("当前下载进度 %f\n", float64(r.NowFileSize*10000/r.FinishFileSize)/100)
//	return
//}
var (
	pixivUserDownUrlRoot string
	pixivPhotoPath       string
	serverLog            *zap.SugaredLogger
	logConfig            logConfigStruct
)

// 神奇的跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

// 配置文件读取器
func init() {
	fmt.Println("---------开始初始化程序---------")
	fmt.Println("---------开始读取配置文件---------")
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("#####初始化失败，没有找到对应的配置文件，程序自动退出#####")
		panic("ERROR:NOT HAVE CONFIG")
	}
	fmt.Println("---------读取配置文件成功---------")
	fmt.Println("---------开始处理配置文件---------")
	pixivUserDownUrlRoot = viper.GetString("basis.url")
	fmt.Println("---------读取网站图片链接成功---------")
	fmt.Println("图片链接:", pixivUserDownUrlRoot)
	pixivPhotoPath = viper.GetString("basis.PhotoPath")
	fmt.Println("---------读取网站存放位置成功---------")
	fmt.Println("图片存放位置:", pixivPhotoPath)
	logConfig.LogPath = viper.GetString("logconfig.LogPath")
	logConfig.MaxSize = viper.GetInt("logconfig.MaxSize")
	logConfig.MaxSaveAge = viper.GetInt("logconfig.MaxSaveAge")
	logConfig.MaxBackup = viper.GetInt("logconfig.MaxBackup")
	fmt.Println("---------读取日志配置成功---------")
	fmt.Println("日志存放位置:", logConfig.LogPath)
	fmt.Println("单个日志最大大小:", logConfig.MaxSize)
	fmt.Println("日志最长存放时间:", logConfig.MaxSaveAge)
	fmt.Println("日志存档最大数:", logConfig.MaxBackup)
	fmt.Println("---------所有配置文件读取成功，开始初始化程序---------")
	fmt.Println("---------正在初始化日志---------")
	InitLogger()
	fmt.Println("---------初始化完成---------")
	fmt.Println("---------欢迎使用鸢月图片代下系统---------")
}

// 日志库实现
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	serverLog = logger.Sugar()
}

//定义时间编码和格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 使用lumberjack实现日志自动分割，默认是最大大小10M
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logConfig.LogPath,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackup,
		MaxAge:     logConfig.MaxSaveAge,
		//暂时不支持用户归档储存
		//TODO 这里以后要支持压缩
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 实现一个Pixiv文件下载器
func downfile(filename, url string) (code, message, downurl string) {

	// TODO 想个办法优化下这里，先用两个URL顶着
	downName := filename
	// 构建一个http请求downtool，携带pixiv的Referer
	filename = pixivPhotoPath + filename
	downrequests := http.Client{}
	pixivRequest, _ := http.NewRequest("GET", url, nil)
	//加入header头
	pixivRequest.Header.Add("referer", "https://www.pixiv.net/")
	pixivRequest.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	//发起请求，取得信息
	response, err := downrequests.Do(pixivRequest)
	if err != nil {
		fmt.Println(err.Error())
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
		userDownUrl := pixivUserDownUrlRoot + downName
		serverLog.Info("文件写入成功，写入文件名" + filename)
		return "200", "文件下载成功，文件大小" + strconv.FormatInt(n/1024, 10) + "KB", userDownUrl

	} else if response.StatusCode == 403 { //处理下因为超载导致的问题
		return "403", "读取PIXIV失败，错误如下:" + err.Error(), ""
	} else {
	}
	return "400", "未知原因下载失败", ""
}

// 实现Pixiv下载的全部接口
func pixivDownFile(context *gin.Context) {
	// 获取用户请求
	pid := context.Query("pid")
	if pid == "" {
		context.JSON(200, gin.H{
			"code":    "400",
			"message": "错误的请求内容",
		})
		context.Abort()
		return

	}
	// 解析用户请求

	code, downUrl := parPixivPid(pid)
	filename := pid + ".png"
	// 检查解析结果是否正常，是否已经被缓存，若已经被缓存，则直接返回缓存结果
	if code == "201" {
		context.JSON(200, gin.H{
			"code":     "200",
			"filename": filename,
			"message":  "命中缓存，已上传CDN",
			"body":     downUrl,
		})
		context.Abort()
		return
	}
	// 检查是否发生错误，发生错误则跳过请求
	if code != "200" {
		context.JSON(200, gin.H{
			"code":    code,
			"message": downUrl,
		})
		context.Abort()
		return
	}

	code, message, downurl := downfile(filename, downUrl)

	if code != "200" {
		context.JSON(200, gin.H{
			"code":    code,
			"message": message,
		})
		context.Abort()
		return
	}
	//待优化
	if code == "200" {
		context.JSON(200, gin.H{
			"code":     "200",
			"filename": filename,
			"message":  "未缓存资源，加入缓存",
			"body":     downurl,
		})
		context.Abort()
		return
	}
	context.Abort()
}

// 新API实现的下载接口，无法判断作品信息，只能下载
func pixivDownDirectFile(context *gin.Context) {
	// 读取ur传递得到参数
	pid := context.Query("pid")
	// 处理链接
	pixivUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages?lang=zh", pid)
	// 构建HTTP请求器
	pixivReClient := http.Client{}
	pixivRequest, err := http.NewRequest("GET", pixivUrl, nil)
	if err != nil {
		context.JSON(200, gin.H{})
		//return "500", fmt.Sprint("错误，构建解析请求失败，错误信息:", err.Error()), *returnInfo
	}
	pixivReClient.Do(pixivRequest)
}

// 实现一个Pixiv链接信息解析（外部版本）
func parsPixivInfo(pid string) (code, message string, responBody parePixivReturn) {
	returnInfo := new(parePixivReturn)
	if pid == "" {
		return "400", fmt.Sprint("错误，无效的请求"), *returnInfo

	}
	parsUrl := fmt.Sprint("https://www.pixiv.net/ajax/illust/", pid)
	getParsClient := http.Client{}
	getParsRequests, err := http.NewRequest("GET", parsUrl, nil)
	if err != nil {
		return "500", fmt.Sprint("错误，构建解析请求失败，错误信息:", err.Error()), *returnInfo
	}
	getParsRequests.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.41 Safari/537.36 Edg/90.0.818.22")
	pareRespon, err := getParsClient.Do(getParsRequests)
	if err != nil {
		return "500", fmt.Sprint("错误，发起请求失败，错误信息:", err.Error()), *returnInfo
	}
	all, err := ioutil.ReadAll(pareRespon.Body)
	if err != nil {
		return "500", fmt.Sprint("错误，解析返回信息失败，错误信息:", err.Error()), *returnInfo
	}
	responJson := parePixivJson{}
	testResponJson := testParePixiv{}
	fmt.Println()
	//测试解析状态
	err = json.Unmarshal(all, &testResponJson)
	if err != nil {
		return "500", fmt.Sprint("错误，基础解析json信息失败,服务端出现问题"), *returnInfo
	}
	if testResponJson.Error == true || testResponJson.Message != "" {
		if testResponJson.Message == "该作品已被删除，或作品ID不存在。" || testResponJson.Message == "該当作品は削除されたか、存在しない作品IDです。" || testResponJson.Message == "リクエストされたページが見つかりませんでした" || testResponJson.Message == "无法找到您所请求的页面" {
			return "404", fmt.Sprint("错误，该作品已经不存在"), *returnInfo
		} else if testResponJson.Message == "不正なリクエストです。" || testResponJson.Message == "不正确的请求。" {
			return "400", fmt.Sprint("错误的请求方式"), *returnInfo
		} else {
			serverLog.Error("发生错误，无法识别的P站报错，报错信息如下:" + testResponJson.Message)
			return "500", fmt.Sprint("未勘测到的错误，请提供PID联系管理员"), *returnInfo
		}

	}

	err = json.Unmarshal(all, &responJson)
	if err != nil {
		return "500", fmt.Sprint("错误，解析json信息失败"), *returnInfo
	}

	returnInfo.Pid = responJson.Body.IllustId
	returnInfo.Name = responJson.Body.Title
	returnInfo.UpdateTime = responJson.Body.UploadDate
	returnInfo.Downurl = pixivDownUrl{
		Mini:     responJson.Body.Urls.Mini,
		Original: responJson.Body.Urls.Original,
	}
	returnInfo.Width = int64(responJson.Body.Width)
	returnInfo.Height = int64(responJson.Body.Height)

	return "200", "OK", *returnInfo
}

// 实现一个Pixiv链接信息解析（内部版本）
func parPixivPid(pid string) (code, url string) {

	dir, err := ioutil.ReadDir(pixivPhotoPath)
	filename := pid + ".png"
	if err != nil {
		return "trst", "test"
	}
	for i := 0; i < len(dir); i++ {
		if (dir[i].Name()) == filename {

			return "201", pixivUserDownUrlRoot + filename
		}
	}
	parsUrl := fmt.Sprint("https://www.pixiv.net/ajax/illust/", pid)
	pareRequestsClient := http.Client{}
	pareRequests, err := http.NewRequest("GET", parsUrl, nil)
	if err != nil {
		return "500", "解析失败，无法构建请求"
	}
	responDo, err := pareRequestsClient.Do(pareRequests)
	if err != nil {
		return "500", "发起请求失败，请检查构建问题"
	}
	respon, err := ioutil.ReadAll(responDo.Body)
	if err != nil {
		return "500", "解析字节流失败"
	}
	responJson := parePixivJson{}
	testResponJson := testParePixiv{}
	err = json.Unmarshal(respon, &testResponJson)
	if testResponJson.Error == true || testResponJson.Message != "" {
		if testResponJson.Message == "该作品已被删除，或作品ID不存在。" || testResponJson.Message == "該当作品は削除されたか、存在しない作品IDです。" || testResponJson.Message == "リクエストされたページが見つかりませんでした" || testResponJson.Message == "无法找到您所请求的页面" {
			return "404", fmt.Sprint("错误，该作品已经不存在")
		} else if testResponJson.Message == "不正なリクエストです。" || testResponJson.Message == "不正确的请求。" {
			return "400", fmt.Sprint("错误的请求方式")
		} else {
			serverLog.Error("发生错误，无法识别的P站报错，报错信息如下:" + testResponJson.Message)
			return "500", fmt.Sprint("未勘测到的错误，请提供PID联系管理员")
		}
	}
	err = json.Unmarshal(respon, &responJson)
	if err != nil {
		fmt.Println(err)
		return "500", "解析json失败"
	}
	if responJson.Body.Urls.Original == "" {
		return "500", "处理失败，解析链接不存在"
	}
	return "200", responJson.Body.Urls.Original
}
func main() {
	server := gin.Default()
	server.Use(Cors())
	server.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "success",
		})
	})
	//实现一个GIN路由组，该组负责接收下载的文件
	pixivServer := server.Group("/pixiv")
	pixivServer.GET("/get/pare/img", func(context *gin.Context) {
		pid := context.Query("pid")
		resultCode, message, respone := parsPixivInfo(pid)
		context.JSON(200, gin.H{
			"code":    resultCode,
			"message": message,
			"body":    respone,
		})
	})
	pixivServer.GET("/get/down/img", func(context *gin.Context) {
		pixivDownFile(context)
	})

	server.Run("0.0.0.0:4560")

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
