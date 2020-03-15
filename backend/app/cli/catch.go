package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type wxBookDetail struct {
	ID                 int         `json:"Id"`
	ColumnID           int         `json:"ColumnId"`
	Code               string      `json:"Code"`
	Name               string      `json:"Name"`
	Author             string      `json:"Author"`
	Tags               string      `json:"Tags"`
	HomeImg            string      `json:"HomeImg"`
	PlayImg            string      `json:"PlayImg"`
	Abstract           string      `json:"Abstract"`
	IsRecommend        int         `json:"IsRecommend"`
	SetRecommendDate   string      `json:"SetRecommendDate"`
	IsHot              int         `json:"IsHot"`
	SetHotDate         string      `json:"SetHotDate"`
	IsSelected         int         `json:"IsSelected"`
	SetSelectedDate    string      `json:"SetSelectedDate"`
	IsSlide            int         `json:"IsSlide"`
	SetSlideDate       string      `json:"SetSlideDate"`
	IsPay              int         `json:"IsPay"`
	PayPrice           float64     `json:"PayPrice"`
	Point              int         `json:"Point"`
	IsComment          int         `json:"IsComment"`
	IsActivity         int         `json:"IsActivity"`
	ClickTimes         int         `json:"ClickTimes"`
	SeoTitle           interface{} `json:"SeoTitle"`
	SeoKeywords        interface{} `json:"SeoKeywords"`
	SeoDescriptions    interface{} `json:"SeoDescriptions"`
	AudioSort          int         `json:"AudioSort"`
	Title              string      `json:"Title"`
	AudioAbstract      string      `json:"AudioAbstract"`
	FileName           string      `json:"FileName"`
	APlayImg           string      `json:"APlayImg"`
	FilePath           string      `json:"FilePath"`
	FileSuffix         string      `json:"FileSuffix"`
	FileSize           int         `json:"FileSize"`
	FileDuration       int         `json:"FileDuration"`
	ListenFileName     interface{} `json:"ListenFileName"`
	ListenFilePath     interface{} `json:"ListenFilePath"`
	ListenFileSize     int         `json:"ListenFileSize"`
	ListenFileSuffix   interface{} `json:"ListenFileSuffix"`
	ListenFileDuration int         `json:"ListenFileDuration"`
	PlayTimes          int         `json:"PlayTimes"`
	AuditStatus        int         `json:"AuditStatus"`
	AuditUserID        int         `json:"AuditUserId"`
	AuditUserRealname  interface{} `json:"AuditUserRealname"`
	AuditDate          interface{} `json:"AuditDate"`
	AuditDesc          interface{} `json:"AuditDesc"`
	AuditLog           interface{} `json:"AuditLog"`
	IsEnabled          int         `json:"IsEnabled"`
	IsDelete           int         `json:"IsDelete"`
	SortCode           int         `json:"SortCode"`
	Remark             interface{} `json:"Remark"`
	CreateDate         string      `json:"CreateDate"`
	AID                int         `json:"AId"`
	AAbstract          string      `json:"AAbstract"`
}

type wxBookList struct {
	ID                 int         `json:"Id"`
	ColumnID           int         `json:"ColumnId"`
	Code               string      `json:"Code"`
	Name               string      `json:"Name"`
	Author             string      `json:"Author"`
	Tags               string      `json:"Tags"`
	HomeImg            string      `json:"HomeImg"`
	PlayImg            string      `json:"PlayImg"`
	Abstract           string      `json:"Abstract"`
	IsRecommend        int         `json:"IsRecommend"`
	SetRecommendDate   interface{} `json:"SetRecommendDate"`
	IsHot              int         `json:"IsHot"`
	SetHotDate         interface{} `json:"SetHotDate"`
	IsSelected         int         `json:"IsSelected"`
	SetSelectedDate    interface{} `json:"SetSelectedDate"`
	IsSlide            int         `json:"IsSlide"`
	SetSlideDate       interface{} `json:"SetSlideDate"`
	IsPay              int         `json:"IsPay"`
	PayPrice           float64     `json:"PayPrice"`
	Point              int         `json:"Point"`
	IsComment          int         `json:"IsComment"`
	IsActivity         int         `json:"IsActivity"`
	ClickTimes         int         `json:"ClickTimes"`
	SeoTitle           string      `json:"SeoTitle"`
	SeoKeywords        string      `json:"SeoKeywords"`
	SeoDescriptions    string      `json:"SeoDescriptions"`
	IsEnabled          int         `json:"IsEnabled"`
	IsDelete           int         `json:"IsDelete"`
	AudioSort          int         `json:"AudioSort"`
	SortCode           int         `json:"SortCode"`
	ColumnSortCode     int         `json:"ColumnSortCode"`
	Remark             interface{} `json:"Remark"`
	CreateDate         string      `json:"CreateDate"`
	CreateUserID       string      `json:"CreateUserId"`
	CreateUserRealname string      `json:"CreateUserRealname"`
	ModifyDate         string      `json:"ModifyDate"`
	ModifyUserID       string      `json:"ModifyUserId"`
	ModifyUserRealname string      `json:"ModifyUserRealname"`
	BookID             int         `json:"BookId"`
	PlayTimes          int         `json:"PlayTimes"`
}

type bookDetail struct {
	Name          string `json:"Name"`
	Title         string `json:"Title"`
	HomeImg       string `json:"HomeImg"`
	AudioAbstract string `json:"AudioAbstract"`
	FileSize      int    `json:"FileSize"`
	FileDuration  int    `json:"FileDuration"`
	CreateDate    string `json:"CreateDate"`
	FilePath      string `json:"FilePath"`
}

type bookList struct {
	ID         int     `json:"Id"`
	Name       string  `json:"Name"`
	Author     string  `json:"Author"`
	HomeImg    string  `json:"HomeImg"`
	Abstract   string  `json:"Abstract"`
	PayPrice   float64 `json:"PayPrice"`
	CreateDate string  `json:"CreateDate"`
	Detail     []bookDetail
}

func getBookData(pageIndex int) (b []bookList, err error) {

	pageSize := 15
	bookListURL := "https://wx.laomassf.com/prointerface/MiniApp/Index.asmx/GetBookList"

	values := map[string]interface{}{"types": "0", "pageIndex": pageIndex, "pageSize": pageSize}
	jsonStr, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	jsonBody := getRequestPost(bookListURL, jsonStr)

	firstData := parser(jsonBody)
	secondData := parser(firstData["d"])

	var wxBookListObj []wxBookList

	if secondData["Data"] == nil {
		return nil, errors.New("data is nil")
	}
	_ = json.Unmarshal([]byte(secondData["Data"].(string)), &wxBookListObj)
	lenWxBookListObj := len(wxBookListObj)
	var bookListObj = make([]bookList, lenWxBookListObj)
	if lenWxBookListObj == 0 {
		return nil, err
	}

	bookIDChan := make(chan int, lenWxBookListObj)
	resultChan := make(chan []bookDetail, lenWxBookListObj)

	for i := 0; i < lenWxBookListObj; i++ {
		bookListObj[i].ID = wxBookListObj[i].ID
		bookListObj[i].Name = wxBookListObj[i].Name
		bookListObj[i].Author = wxBookListObj[i].Author
		bookListObj[i].HomeImg = urlPathFormat(wxBookListObj[i].HomeImg)
		bookListObj[i].Abstract = wxBookListObj[i].Abstract
		bookListObj[i].PayPrice = wxBookListObj[i].PayPrice
		bookListObj[i].CreateDate = createDateFormat(wxBookListObj[i].CreateDate)
		bookListObj[i].Detail = nil

	}
	go worker(bookIDChan, resultChan)

	for i := 0; i < len(bookListObj); i++ {
		bookIDChan <- bookListObj[i].ID
	}
	close(bookIDChan)

	for i := 0; i < len(bookListObj); i++ {
		bookListObj[i].Detail = <-resultChan
	}

	return bookListObj, nil

}

func worker(bookIDChan <-chan int, results chan<- []bookDetail) {
	for bookID := range bookIDChan {
		results <- getBookDetail(bookID)
	}
}

func getRequestPost(urlStr string, jsonStr []byte) string {

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "wx.laomassf.com")
	req.Header.Set("Referer", "https://servicewechat.com/wx1f8180176500f209/25/page-frame.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.10(0x17000a21) NetType/WIFI Language/zh_CN")
	// proxyURL, err := url.Parse("http://127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	// tr := &http.Transport{
	// 	Proxy: http.ProxyURL(proxyURL),
	// }
	// client := &http.Client{Transport: tr}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func getBookDetail(bookID int) []bookDetail {
	fmt.Printf("%d now ", bookID)

	bookDetailURL := "https://wx.laomassf.com/prointerface/MiniApp/Index.asmx/GetAudioList"
	values := map[string]interface{}{"bookId": bookID}
	jsonStr, _ := json.Marshal(values)

	jsonBody := getRequestPost(bookDetailURL, jsonStr)

	firstData := parser(jsonBody)
	secondData := parser(firstData["d"])
	var wxBooksObj []wxBookDetail
	_ = json.Unmarshal([]byte(secondData["Data"].(string)), &wxBooksObj)
	var bookDetailObj = make([]bookDetail, len(wxBooksObj))
	for i := 0; i < len(wxBooksObj); i++ {
		bookDetailObj[i].Name = wxBooksObj[i].Name
		bookDetailObj[i].Title = wxBooksObj[i].Title
		bookDetailObj[i].HomeImg = urlPathFormat(wxBooksObj[i].HomeImg)
		bookDetailObj[i].AudioAbstract = wxBooksObj[i].AudioAbstract
		bookDetailObj[i].FileSize = wxBooksObj[i].FileSize
		bookDetailObj[i].FileDuration = wxBooksObj[i].FileDuration
		bookDetailObj[i].CreateDate = createDateFormat(wxBooksObj[i].CreateDate)
		bookDetailObj[i].FilePath = urlPathFormat(wxBooksObj[i].FilePath)
	}
	return bookDetailObj
}

func urlPathFormat(urlPath string) string {
	return "https://wx.laomassf.com" + urlPath
}

func example() {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	page := 1
	for {
		var bookObj []bookList
		bookObj, err := getBookData(page)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		for _, val := range bookObj {
			mapBookList := make(map[string]interface{})

			jsonBookList, _ := json.Marshal(val)
			json.Unmarshal(jsonBookList, &mapBookList)

			jsonDetail, _ := json.Marshal(mapBookList["Detail"])
			delete(mapBookList, "Detail")

			mapBookList["Detail"] = string(jsonDetail)

			err := client.HMSet(strconv.Itoa(val.ID), mapBookList).Err()
			if err != nil {
				panic(err)
			}
		}
		fmt.Println(page)
		page++
	}

}

func main() {
	example()
}

func parser(data interface{}) map[string]interface{} {
	var i interface{}
	json.Unmarshal([]byte(data.(string)), &i)
	jData, _ := i.(map[string]interface{})
	return jData
}

func createDateFormat(createDate string) string {
	i, _ := strconv.ParseInt(createDate[6:len(createDate)-5], 10, 64)
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02 15:04:05")
}
