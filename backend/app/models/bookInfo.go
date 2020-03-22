package models

import (
	"encoding/json"
	"strconv"

	"app/dao"

	"github.com/go-redis/redis"
)

type BookInfo struct {
}

type bookDetail struct {
	Title         string `json:"Title"`
	AudioAbstract string `json:"AudioAbstract"`
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

func (d *BookInfo) FindBookDetailById(id int) ([]bookDetail, error) {
	client := dao.RedisClient()
	var err error
	client.Do("select", 0)
	item := client.HGetAll(strconv.Itoa(id))
	var bookDetailList []bookDetail
	err = json.Unmarshal([]byte(item.Val()["Detail"]), &bookDetailList)
	var bookDetailObj = make([]bookDetail, len(bookDetailList))
	for i := 0; i < len(bookDetailList); i++ {
		bookDetailObj[i].Title = bookDetailList[i].Title
		bookDetailObj[i].AudioAbstract = bookDetailList[i].AudioAbstract
		bookDetailObj[i].CreateDate = bookDetailList[i].CreateDate
		bookDetailObj[i].FilePath = bookDetailList[i].FilePath
	}
	return bookDetailObj, err
}

func (d *BookInfo) FindBookListByPage() (map[string]bookList, error) {
	client := dao.RedisClient()
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
	return bookListObj, err
}
