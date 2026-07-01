package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
)

var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "bitedash",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

var HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "bitedash",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

var GRPCRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "bitedash",
		Subsystem: "grpc",
		Name:      "requests_total",
		Help:      "Total number of gRPC requests.",
	},
	[]string{"method", "code"},
)

var GRPCRequestDurationSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "bitedash",
		Subsystem: "grpc",
		Name:      "request_duration_seconds",
		Help:      "gRPC request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "code"},
)

func Register() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDurationSeconds,
		GRPCRequestsTotal,
		GRPCRequestDurationSeconds,
	)
}

func ObserveHTTPRequest(c *gin.Context, duration time.Duration) {
	method := c.Request.Method
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	status := strconv.Itoa(c.Writer.Status())

	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDurationSeconds.WithLabelValues(method, path, status).Observe(duration.Seconds())
}

func ObserveGRPCRequest(method string, code codes.Code, duration time.Duration) {
	codeStr := code.String()

	GRPCRequestsTotal.WithLabelValues(method, codeStr).Inc()
	GRPCRequestDurationSeconds.WithLabelValues(method, codeStr).Observe(duration.Seconds())
}
