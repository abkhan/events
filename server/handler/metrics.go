package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsHandler(version string) gin.HandlerFunc {
	infoGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name:      "info",
		Help:      "Service info",
		Namespace: "fixme_ns",
		ConstLabels: map[string]string{
			"version": version,
		},
	})

	infoGauge.Inc()

	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
