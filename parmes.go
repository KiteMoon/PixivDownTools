package main

type testParePixiv struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
type parePixivJson struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    struct {
		IllustId      string `json:"illustId"`
		IllustTitle   string `json:"illustTitle"`
		IllustComment string `json:"illustComment"`
		Id            string `json:"id"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		IllustType    int    `json:"illustType"`
		CreateDate    string `json:"createDate"`
		UploadDate    string `json:"uploadDate"`
		Restrict      int    `json:"restrict"`
		XRestrict     int    `json:"xRestrict"`
		Sl            int    `json:"sl"`
		Urls          struct {
			Mini     string `json:"mini"`
			Thumb    string `json:"thumb"`
			Small    string `json:"small"`
			Regular  string `json:"regular"`
			Original string `json:"original"`
		} `json:"urls"`
		Tags                    interface{}   `json:"tags"`
		Alt                     string        `json:"alt"`
		StorableTags            []string      `json:"storableTags"`
		UserId                  string        `json:"userId"`
		UserName                string        `json:"userName"`
		UserAccount             string        `json:"userAccount"`
		UserIllusts             interface{}   `json:"userIllusts"`
		LikeData                bool          `json:"likeData"`
		Width                   int           `json:"width"`
		Height                  int           `json:"height"`
		PageCount               int           `json:"pageCount"`
		BookmarkCount           int           `json:"bookmarkCount"`
		LikeCount               int           `json:"likeCount"`
		CommentCount            int           `json:"commentCount"`
		ResponseCount           int           `json:"responseCount"`
		ViewCount               int           `json:"viewCount"`
		BookStyle               int           `json:"bookStyle"`
		IsHowto                 bool          `json:"isHowto"`
		IsOriginal              bool          `json:"isOriginal"`
		ImageResponseOutData    []interface{} `json:"imageResponseOutData"`
		ImageResponseData       []interface{} `json:"imageResponseData"`
		ImageResponseCount      int           `json:"imageResponseCount"`
		PollData                interface{}   `json:"pollData"`
		SeriesNavData           interface{}   `json:"seriesNavData"`
		DescriptionBoothId      interface{}   `json:"descriptionBoothId"`
		DescriptionYoutubeId    interface{}   `json:"descriptionYoutubeId"`
		ComicPromotion          interface{}   `json:"comicPromotion"`
		FanboxPromotion         interface{}   `json:"fanboxPromotion"`
		ContestBanners          []interface{} `json:"contestBanners"`
		IsBookmarkable          bool          `json:"isBookmarkable"`
		BookmarkData            interface{}   `json:"bookmarkData"`
		ContestData             interface{}   `json:"contestData"`
		ZoneConfig              interface{}   `json:"zoneConfig"`
		ExtraData               interface{}   `json:"extraData"`
		TitleCaptionTranslation interface{}   `json:"titleCaptionTranslation"`
		IsUnlisted              bool          `json:"isUnlisted"`
		Request                 interface{}   `json:"request"`
		CommentOff              int           `json:"commentOff"`
		NoLoginData             interface{}   `json:"noLoginData"`
	} `json:"body"`
}

type parePixivReturn struct {
	Pid        string
	Name       string
	UpdateTime string
	Downurl    pixivDownUrl
	Width      int64
	Height     int64
}
type pixivDownUrl struct {
	Mini     string
	Original string
}
type PixivDirectPareBodyListUrl struct {
	ThumbMini string `json:"thumb_mini"`
	Small     string `json:"small"`
	Regular   string `json:"regular"`
	Original  string `json:"original"`
}
type PixivDirectPareBodyList struct {
	Urls   PixivDirectPareBodyListUrl `json:"urls"`
	Width  int                        `json:"width"`
	Height int                        `json:"height"`
}
type PixivDirectPare struct {
	Error   bool                      `json:"error"`
	Message string                    `json:"message"`
	Body    []PixivDirectPareBodyList `json:"body"`
}

// 多PAGE下载器
type PidDownToolsMore struct {
	Name string
	Num  int `json:"num"`
	Body []PidDownToolsMoreBody
}

// 多PAGE下载器Body部分（复用）
type PidDownToolsMoreBody struct {
	Url    string
	Width  int
	Height int
}

// 多PAGE下载器传回端
type PidDownToolsRespon struct {
	Num        int `json:"num"`
	SuccessNum int `json:"success"`
	ErrorNum   int `json:"errorNum"`
	Body       []PidDownToolsMoreBody
}

// 日志库配置文件实现
type logConfigStruct struct {
	LogPath    string
	MaxSize    int
	MaxSaveAge int
	MaxBackup  int
}
