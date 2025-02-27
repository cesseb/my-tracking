package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Car represents the structure for car details
type Car struct {
	Model     string  `json:"model"`
	Price     float64 `json:"price"`
	RegNumber string  `json:"reg_number"`
}

// A slice to hold the car details
var cars = []Car{
	{"Toyota Camry", 24000, "ABC123"},
	{"Honda Accord", 23000, "XYZ456"},
	{"Tesla Model 3", 35000, "TESLA789"},
}

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, activeConnections)
}

// Middleware to track Prometheus metrics
func prometheusMiddleware(c *gin.Context) {
	path := c.Request.URL.Path

	// begin timer to measure the requests duration
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(path))

	// increment total request counter
	httpRequestsTotal.WithLabelValues(path).Inc()

	// increment number of active connections
	activeConnections.Inc()

	// complete processing request
	c.Next()

	// record request duration (post processing)
	timer.ObserveDuration()

	// decrement total number of active connections (post processing)
	activeConnections.Dec()
}

func main() {
	router := gin.Default()

	// Middleware to collect metrics
	router.Use(prometheusMiddleware)

	// Endpoint to get car details
	router.GET("/cars", func(c *gin.Context) {
		c.JSON(http.StatusOK, cars)
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start the server
	router.Run(":8081")
}
