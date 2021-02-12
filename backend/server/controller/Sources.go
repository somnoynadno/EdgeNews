package controller

import (
	"EdgeNews/backend/db/dao"
	u "EdgeNews/backend/utils"
	"encoding/json"
	"net/http"
)

var GetAllSources = func(w http.ResponseWriter, r *http.Request) {
	sources, err := dao.GetAllSources()
	if err != nil {
		u.HandleInternalError(w, err)
	} else {
		res, _ := json.Marshal(sources)
		u.RespondJSON(w, res)
	}
}