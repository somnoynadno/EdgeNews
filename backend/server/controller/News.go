package controller

import (
	"EdgeNews/backend/db/dao"
	u "EdgeNews/backend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var GetLastNews = func(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["amount"]
	if !ok || len(keys[0]) < 1 {
		u.HandleBadRequest(w, errors.New("parameter 'amount' is missing"))
		return
	}

	amount, err := strconv.Atoi(keys[0])
	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	news, err := dao.GetLastNews(amount)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		res, _ := json.Marshal(news)
		u.RespondJSON(w, res)
	}
}
