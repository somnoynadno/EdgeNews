package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Responses responses
	Scrapings scrapings
	WS        websocket
}

type responses struct {
	OK            prometheus.Counter
	BadRequest    prometheus.Counter
	InternalError prometheus.Counter
}

type scrapings struct {
	Done   prometheus.CounterVec
	Failed prometheus.CounterVec
}

type websocket struct {
	ConnectionsActive prometheus.Gauge
	BroadcastMessages prometheus.Counter
}

var metrics *Metrics

var (
	labels = []string{"scrapper"}

	promResponseOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "edge_api_http_response_200",
		Help: "Total number of 200 requests",
	})
	promResponseBadRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "edge_api_http_response_400",
		Help: "Total number of 400 requests",
	})
	promResponseInternalError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "edge_api_http_response_500",
		Help: "Total number of 500 requests",
	})

	promScrapingsDone = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "edge_api_http_scrapings_done",
		Help: "Number of processed scrapings",
	}, labels)
	promScrapingsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "edge_api_http_scrapings_failed",
		Help: "Number of failed scrapings",
	}, labels)

	promHubConnectionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "edge_ws_active_connections",
		Help: "Number of active connections on websocket",
	})
	promHubBroadcastMessages = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "edge_ws_broadcast_messages",
		Help: "Number of sent websocket messages",
	})
)

func init() {
	prometheus.MustRegister(promHubBroadcastMessages)
	prometheus.MustRegister(promHubConnectionsActive)

	metrics = &Metrics{
		Responses: responses{
			OK:            promResponseOK,
			BadRequest:    promResponseBadRequest,
			InternalError: promResponseInternalError,
		},
		Scrapings: scrapings{
			Done:   *promScrapingsDone,
			Failed: *promScrapingsFailed,
		},
		WS: websocket{
			ConnectionsActive: promHubConnectionsActive,
			BroadcastMessages: promHubBroadcastMessages,
		},
	}
}

func GetMetrics() *Metrics {
	return metrics
}
