package controller

import (
	"EdgeNews/backend/db/dao"
	u "EdgeNews/backend/utils"
	"encoding/json"
	"net/http"
)

var GetActiveTextStreams = func(w http.ResponseWriter, r *http.Request) {
	textStreams, err := dao.GetActiveTextStreams()
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		res, _ := json.Marshal(textStreams)
		u.RespondJSON(w, res)
	}
}
