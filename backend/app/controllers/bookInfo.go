package controllers

import (
	"app/config"
	"app/helper"
	"app/models"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	model = models.BookInfo{}
)

func Home(w http.ResponseWriter, r *http.Request) {
	helper.ResponseWithJson(w, 0, "welcome")
}

func BookList(w http.ResponseWriter, r *http.Request) {
	result, err := model.FindBookListByPage()
	if err != nil {
		helper.ResponseWithJson(w, config.ErrServiceBusy, err.Error())
		return
	}
	helper.ResponseWithJson(w, 0, result)
}

func BookDetail(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	bookID := -1
	if val, ok := pathParams["id"]; ok {
		var err error
		bookID, err = strconv.Atoi(val)
		if err != nil {
			helper.ResponseWithJson(w, config.ErrBookId, err.Error())
			return
		}
	}
	result, err := model.FindBookDetailById(bookID)
	if err != nil {
		helper.ResponseWithJson(w, config.ErrBookIdOverFlow, err.Error())
		return
	}
	helper.ResponseWithJson(w, 0, result)

}
