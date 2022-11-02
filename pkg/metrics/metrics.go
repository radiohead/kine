package metrics

import (
	"time"

	"github.com/k3s-io/kine/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	ResultSuccess = "success"
	ResultError   = "error"
)

var (
	SQLTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kine_sql_total",
		Help: "Total number of SQL operations",
	}, []string{"error_code"})

	SQLTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "kine_sql_time_seconds",
		Help: "Length of time per SQL operation",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
			1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30},
	}, []string{"error_code"})

	CompactTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kine_compact_total",
		Help: "Total number of compactions",
	}, []string{"result"})

	SQLWatchGoroutineCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kine_sql_active_watch_goroutines",
		Help: "Number of active WATCH goroutines",
	})

	SQLCompactionTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "kine_sql_compaction_latency_seconds",
		Help:    "Histogram measuring the latency of database compaction operations",
		Buckets: prometheus.ExponentialBuckets(0.001, 4, 8),
	}, []string{"result"})

	SQLTTLCacheSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kine_sql_ttl_deletion_cache_size",
		Help: "Gauge measuring the size of the cache used by TTL expiration",
	})

	SQLTTLDeletionTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "kine_sql_ttl_deletion_latency_seconds",
		Help:    "Histogram measuring the latency of deletions caused by TTL expiration",
		Buckets: prometheus.ExponentialBuckets(0.001, 4, 8),
	}, []string{"result"})
)

var (
	// SlowSQLThreshold is a duration which SQL executed longer than will be logged.
	// This can be directly modified to override the default value when kine is used as a library.
	SlowSQLThreshold = time.Second
)

func ObserveSQL(start time.Time, errCode string, sql util.Stripped, args ...interface{}) {
	SQLTotal.WithLabelValues(errCode).Inc()
	duration := time.Since(start)
	SQLTime.WithLabelValues(errCode).Observe(duration.Seconds())
	if SlowSQLThreshold > 0 && duration >= SlowSQLThreshold {
		logrus.Infof("Slow SQL (started: %v) (total time: %v): %s : %v", start, duration, sql, args)
	}
}

// ObserveSQLCompaction
func ObserveSQLCompaction(start time.Time, err error) {
	res := ResultSuccess
	if err != nil {
		res = ResultError
	}

	SQLCompactionTime.WithLabelValues(res).Observe(time.Since(start).Seconds())
}
