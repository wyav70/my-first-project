package main

import (
	"encoding/json"
	"log"
	"log/slog"

	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"net/http"
	"strings"
	"time"
)

type MovieRawInfo struct {
	Context         string          `json:"@context"`
	Name            string          `json:"name"`
	URL             string          `json:"url"`
	Image           string          `json:"image"`
	Director        []Director      `json:"director"`
	Author          []interface{}   `json:"author"`
	Actor           []interface{}   `json:"actor"`
	DatePublished   string          `json:"datePublished"`
	Genre           []string        `json:"gener"`
	Duration        string          `json:"duration"`
	Description     string          `json:"decription"`
	Type            string          `json:"@type"`
	AggregateRating AggregateRating `json:"aggregateRating"`
}
type Director struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
	Name string `json:"name"`
}
type AggregateRating struct {
	Type        string `json:"@type"`
	RatingCount string `json:"ratingCount"`
	BestRating  string `json:"bestRating"`
	WorstRating string `json:"worestRating"`
	RatingValue string `json:"ratingValue"`
}
type MovieInfo struct {
	Name          string `json:"name"`
	PublishedDate string `json:"publisheDate"`
	RatingValue   string `json:"rating_value"`
	Desc          string `json:"desc"`
}

// 需求：从豆瓣中获取“出走的决心”的信息
// 1.假设已经知道了电影的id，从豆瓣网中找到该电影的信息
// 2.从找到的信息中定位到想要的内容，电影名字；上映日期；评分；简介
// 3.将信息中定位到的内容制作成表格
func GetDouBanMovieInfo123(url string) string {
	// 定义了一个url 是一个豆瓣的电影信息的地址

	// go 变量的使用是 ,先声明 后使用 ,不声明 无法使用, 会报错的,
	// 声明了却没有使用,也会报错

	// 怎么去获取url内容
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	// 要进行错误的处理
	if err != nil {
		return err.Error()
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	// 要进行异常情况的处理  418 I'm a teapot
	// TODO 要解决这个418问题 需要在headers 里面添加 User-Agent信息
	if res.StatusCode != 200 {
		return "这个页面报错了"
	}
	// 读取资源数据 body: []byte  nil
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)

}

// func 函数的关键字
//
//	GetMovieName 定义的函数名 (infoData string)
//	入参 infoData 定义的参数的名字 (形参)  string 参数的类型 字符串类型
//
// 出参  string 类型
func GetMovieName(infoData string) string {
	newRawList := strings.Split(infoData, `<meta property="og:title" content="`) // 字符串的切割方法 第一个参数是要切割的数据,第二个参数是 要切割的分割标识
	// fmt.Println(newRawList[1])
	newRawStr1 := newRawList[1] // 取第二个参数, 去掉了前面的内容的数据,下面要去掉后面内容
	// fmt.Println(newRawStr1)
	newRawList2 := strings.Split(newRawStr1, "\n")
	newRawStr2 := newRawList2[0]
	// fmt.Println(newRawStr2)

	// 是因为只剩下了</script>  ,所以我们要去掉这个 ,所以要替换成空字符串
	newRaw := strings.ReplaceAll(newRawStr2, `" />`, "")
	// 去除结尾的 "</script>"，替换为空字符串
	// fmt.Println(newRaw)

	// newRaw = strings.ReplaceAll(newRaw, "</script>", "")
	// 去掉前边和后边的空格
	name := strings.TrimSpace(newRaw)
	// fmt.Println(name)
	return name
}

func ProcessDouBanInfoData123(infoData string) MovieRawInfo {
	// 获取电影名称
	movieTitle := GetMovieName(infoData)
	// fmt.Println(movieTitle )
	// 获取电影信息

	preFlag := `<script type="application/ld+json">`
	// 要替换下 原始数据中的电影名称      使用字符串占位符  %s 字符串      %d  数字
	backFlag := fmt.Sprintf(`<meta property="og:title" content="%s" />`, movieTitle)
	// fmt.Println(backFlag )
	// 第一步要去掉 不需要的前半部分内容
	newRawList := strings.Split(infoData, preFlag)
	newRawStr1 := newRawList[1]
	newRawList2 := strings.Split(newRawStr1, backFlag)
	newRawStr2 := newRawList2[0]
	// fmt.Println(newRawStr2 )
	newRaw := strings.ReplaceAll(newRawStr2, `</script>`, "")
	movieInfo := strings.TrimSpace(newRaw)
	// fmt.Println(movieInfo)

	//要去掉前边和后边的空格
	//通过反序列化
	var rawInfo MovieRawInfo
	// // 通过反序列化去把newRaw赋值给rawInfo
	err := json.Unmarshal([]byte(movieInfo), &rawInfo)
	if err != nil {
		log.Fatalf("反序列化错误了，错误内容为：%v", err.Error())
		return rawInfo
	}
	// //返回这个对象，给下一个函数使用，输出对象
	return rawInfo

}
func main() {
	movieIDs := []string{"25824686", "26747919", "36211169", "35087675"}
	//fmt.Println("序号 电影名称 上映时间 评分")
	//创建excel文件
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("movies")
	if err != nil {
		slog.Error("add sheet is failed! ,err:", err)
		return
	}
	//写入表头
	row := sheet.AddRow()
	row.AddCell().Value = "序号"
	row.AddCell().Value = "电影名称"
	row.AddCell().Value = "上映时间"
	row.AddCell().Value = "评分"
	for index, id := range movieIDs {
		url := fmt.Sprintf("https://movie.douban.com/subject/%s/", id)
		result := GetDouBanMovieInfo123(url)
		movieInfo := ProcessDouBanInfoData123(result)

		name := movieInfo.Name
		datePublished := movieInfo.DatePublished
		ratingValue := movieInfo.AggregateRating.RatingValue
		//fmt.Println(index+1, name, datePublished, ratingValue)
		row := sheet.AddRow()
		row.AddCell().Value = fmt.Sprintf("%d", index+1)
		row.AddCell().Value = name
		row.AddCell().Value = datePublished
		row.AddCell().Value = ratingValue
		fmt.Println(name)
	}
	err = file.Save("movies.xlsx")
	if err != nil {
		slog.Error("无法保存文件 %v", err)
		return
	}
	fmt.Println("数据已写入 movies.xlsx 文件.")

}
