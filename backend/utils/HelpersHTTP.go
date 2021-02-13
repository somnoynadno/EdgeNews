package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleOK(w http.ResponseWriter, data map[string]interface{}) {
	log.Info("OK")
	metrics.Responses.OK.Inc()
	Respond(w, data)
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	_ = json.NewEncoder(w).Encode(data)
}

func RespondJSON(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	metrics.Responses.OK.Inc()
	_, _ = w.Write(data)
}

var HandleOptions = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
