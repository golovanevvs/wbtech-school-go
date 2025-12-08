package pkgPrometheus

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"status_code", "method", "path"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	MemoryAllocBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Number of bytes currently allocated",
		},
		[]string{"type"},
	)

	MemorySysBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_memory_sys_bytes",
			Help: "Number of bytes obtained from system",
		},
		[]string{"type"},
	)

	GcCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_gc_count_total",
			Help: "Total number of garbage collections",
		},
		[]string{"gc_type"},
	)

	GcPauseSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_gc_pause_seconds",
			Help:    "Duration of GC pauses",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 20), // от 1ms до ~500s
		},
		[]string{"gc_type"},
	)

	GcLastTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_gc_last_duration_seconds",
			Help: "Duration of last GC",
		},
		[]string{"gc_type"},
	)
)

type gcStatsCollector struct {
	mu        sync.Mutex
	lastNumGC uint32
	lastPause time.Duration
}

var gcCollector = &gcStatsCollector{}

func Init() {
	prometheus.MustRegister(
		RequestCount,
		RequestDuration,
		MemoryAllocBytes,
		MemorySysBytes,
		GcCount,
		GcPauseSeconds,
		GcLastTime,
	)

	go collectMemoryMetrics()
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func collectMemoryMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updateMemoryMetrics()
	}
}

func updateMemoryMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	MemoryAllocBytes.WithLabelValues("heap").Set(float64(memStats.HeapAlloc))
	MemoryAllocBytes.WithLabelValues("stack").Set(float64(memStats.StackInuse))
	MemoryAllocBytes.WithLabelValues("other").Set(float64(memStats.MSpanInuse + memStats.MCacheInuse))

	MemorySysBytes.WithLabelValues("heap").Set(float64(memStats.HeapSys))
	MemorySysBytes.WithLabelValues("stack").Set(float64(memStats.StackSys))
	MemorySysBytes.WithLabelValues("other").Set(float64(memStats.MSpanSys + memStats.MCacheSys))

	gcCollector.mu.Lock()
	defer gcCollector.mu.Unlock()

	if memStats.NumGC > gcCollector.lastNumGC {
		gcCount := memStats.NumGC - gcCollector.lastNumGC
		GcCount.WithLabelValues("total").Add(float64(gcCount))

		if memStats.PauseNs[0] != 0 {
			pauseDuration := time.Duration(memStats.PauseNs[0])
			GcPauseSeconds.WithLabelValues("last").Observe(pauseDuration.Seconds())
			GcLastTime.WithLabelValues("last").Set(pauseDuration.Seconds())
		}

		gcCollector.lastNumGC = memStats.NumGC
	}
}
