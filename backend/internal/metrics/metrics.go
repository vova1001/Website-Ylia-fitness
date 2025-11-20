package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var ReqCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_req_total",
	Help: " Total number of HTTP requests",
},
	[]string{"path", "method", "status"})

var ReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_req_duration_sec",
	Help:    "HTTP req duration in seconds",
	Buckets: prometheus.DefBuckets,
}, []string{"path"})

func MetricsInit() {
	prometheus.MustRegister(ReqCount, ReqDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ctx.Writer.Status())

		ReqCount.WithLabelValues(ctx.FullPath(), ctx.Request.Method, status).Inc()
		ReqDuration.WithLabelValues(ctx.FullPath()).Observe(duration)
	}
}
