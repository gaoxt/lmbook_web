package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/app/singletonRedis"
	"github.com/go-redis/redis"

	"github.com/gorilla/mux"
)

type BaseResponse struct {
	Code    int
	Message string
}

type responseBookDetail struct {
	BaseResponse
	Data []bookDetail
}

type responseBookList struct {
	BaseResponse
	// Data []bookList
	Data map[string]bookList
}

type bookDetail struct {
	Title         string `json:"Title"`
	AudioAbstract string `json:"AudioAbstract"`
	FileSize      string `json:"FileSize"`
	FileDuration  string `json:"FileDuration"`
	CreateDate    string `json:"CreateDate"`
	FilePath      string `json:"FilePath"`
}

type bookList struct {
	ID         string `json:"Id"`
	Name       string `json:"Name"`
	Author     string `json:"Author"`
	HomeImg    string `json:"HomeImg"`
	Abstract   string `json:"Abstract"`
	PayPrice   string `json:"PayPrice"`
	CreateDate string `json:"CreateDate"`
}

type BL struct {
	Key   string
	Value bookList
}
type BookListSort []BL

func (b BookListSort) Len() int      { return len(b) }
func (b BookListSort) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b BookListSort) Less(i, j int) bool {
	return b[i].Value.ID < b[j].Value.ID
}

func getBookDetail(bookID int) *responseBookDetail {
	client := singletonRedis.GetRedis()
	client.Do("select", 0)
	item := client.HGetAll(strconv.Itoa(bookID))
	var bookDetailList []bookDetail
	_ = json.Unmarshal([]byte(item.Val()["Detail"]), &bookDetailList)
	var bookDetailObj = make([]bookDetail, len(bookDetailList))
	for i := 0; i < len(bookDetailList); i++ {
		bookDetailObj[i].Title = bookDetailList[i].Title
		bookDetailObj[i].AudioAbstract = bookDetailList[i].AudioAbstract
		bookDetailObj[i].FileSize = bookDetailList[i].FileSize
		bookDetailObj[i].FileDuration = bookDetailList[i].FileDuration
		bookDetailObj[i].CreateDate = bookDetailList[i].CreateDate
		bookDetailObj[i].FilePath = bookDetailList[i].FilePath
	}
	res := &responseBookDetail{}
	res.Code = 0
	res.Message = "success"
	res.Data = bookDetailObj
	return res
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"code":0,"message":"welcome"}`))
}

func apiBookList(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	page := 1
	var err error
	if val, ok := pathParams["page"]; ok {
		page, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code":1,"message":"error params"}`))
			return
		}
	}
	var jsonObj []byte
	res := getBookList(page)
	jsonObj, _ = json.Marshal(res)
	w.Write([]byte(jsonObj))
}

func getBookList(pageIndex int) *responseBookList {
	client := singletonRedis.GetRedis()
	var keys []string
	var err error
	keys, _, err = client.Scan(0, "", 1000).Result()
	if err != nil {
		panic(err)
	}
	details := map[string]*redis.SliceCmd{}
	pipe := client.Pipeline()
	for _, val := range keys {
		details[val] = pipe.HMGet(val, "Id", "Name", "Author", "HomeImg", "Abstract", "PayPrice", "CreateDate")
	}
	pipe.Exec()

	bookListObj := map[string]bookList{}
	for i, item := range details {
		bookListObj[i] = bookList{
			ID:         item.Val()[0].(string),
			Name:       item.Val()[1].(string),
			Author:     item.Val()[2].(string),
			HomeImg:    item.Val()[3].(string),
			Abstract:   item.Val()[4].(string),
			PayPrice:   item.Val()[5].(string),
			CreateDate: item.Val()[6].(string),
		}
	}
	// var bookListSlice = make([]bookList, len(bookListObj))
	// i := 0
	// for _, val := range bookListObj {
	// 	bookListSlice[i] = val
	// 	i++
	// }
	res := &responseBookList{}
	res.Code = 0
	res.Message = "success"
	res.Data = bookListObj
	return res

}

func getRequestPost(url string, jsonStr []byte) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "wx.laomassf.com")
	req.Header.Set("Referer", "https://servicewechat.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.10(0x17000a21) NetType/WIFI Language/zh_CN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func parser(data interface{}) map[string]interface{} {
	var i interface{}
	json.Unmarshal([]byte(data.(string)), &i)
	jData, _ := i.(map[string]interface{})
	return jData
}

func apiBookDetail(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var jsonObj []byte
	bookID := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		bookID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code":1,"message":"error params"}`))
			return
		}
	}
	if bookID <= 0 {
		var bookDetailObj = make([]bookDetail, 0)
		res := responseBookDetail{}
		res.Code = 2
		res.Message = "error bookID"
		res.Data = bookDetailObj
		jsonObj, _ = json.Marshal(res)
		w.Write([]byte(jsonObj))
	} else {
		res := getBookDetail(bookID)
		jsonObj, _ = json.Marshal(res)
		w.Write([]byte(jsonObj))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", get).Methods(http.MethodGet)
	r.HandleFunc("/list/{page}", apiBookList).Methods(http.MethodGet)
	r.HandleFunc("/detail/{id}", apiBookDetail).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", r))
}

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}
