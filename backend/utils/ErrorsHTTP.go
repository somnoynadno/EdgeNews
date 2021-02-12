package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 400
func HandleBadRequest(w http.ResponseWriter, err error) {
	log.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	metrics.Responses.BadRequest.Inc()
	Respond(w, Message(false, err.Error()))
}

// 500
func HandleInternalError(w http.ResponseWriter, err error) {
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	metrics.Responses.InternalError.Inc()
	Respond(w, Message(false, err.Error()))
}

