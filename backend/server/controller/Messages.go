package controller

import (
	"EdgeNews/backend/db/dao"
	u "EdgeNews/backend/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetMessagesByTextStreamID = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	messages, err := dao.GetMessagesByTextStreamID(id)
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		res, _ := json.Marshal(messages)
		u.RespondJSON(w, res)
	}
}
