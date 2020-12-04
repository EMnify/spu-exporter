package main

import (
	"os"
	"strings"

	"github.com/EMnify/spu-exporter/pkg/collector"
	"github.com/EMnify/spu-exporter/pkg/config"
	"github.com/EMnify/spu-exporter/pkg/prom"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	cfg := config.ReadConfig("configs/config.yml")
	logger := setupLogging(cfg)
	d := collector.NewSpuMetricsDaemon(cfg, logger)
	// Currently config is not read correctly
	trans, err := d.ExecuteScrape()
	if err != nil {
		return
	}
	// prom format
	reg := prom.CreateMetricLines(trans)
	err = prom.WriteToFile(prometheus.Gatherers{reg}, cfg.Prometheus.Outfile)
	if err != nil {
		level.Error(logger).Log("Failed to write results to file: %s", err)
	}
}

func setupLogging(cfg *config.AppConfig) log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	switch strings.ToLower(cfg.LogLevel) {
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	return logger
}
