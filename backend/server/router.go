package server

import (
	"EdgeNews/backend/server/controller"
	"EdgeNews/backend/server/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InitRouter() {
	newsHub := GetNewsHub()
	textStreamHub := GetTextStreamHub()

	http.HandleFunc("/ws/news", func(w http.ResponseWriter, r *http.Request) {
		HandleConnection(newsHub, w, r)
	})

	http.HandleFunc("/ws/streams", func(w http.ResponseWriter, r *http.Request) {
		HandleConnection(textStreamHub, w, r)
	})

	r := mux.NewRouter()

	r.HandleFunc("/api/sources", controller.GetAllSources).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/active_streams", controller.GetActiveTextStreams).Methods(http.MethodGet, http.MethodOptions)
	r.Handle("/metrics", promhttp.Handler())
	http.Handle("/", r)

	r.Use(middleware.CORS)    // enable CORS headers
	r.Use(middleware.LogPath) // log IP, path and method

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

