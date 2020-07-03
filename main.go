package main

import (
	"os"
	"strings"

	"github.com/EMnify/spu-exporter/pkg/config"
	"github.com/EMnify/spu-exporter/pkg/scraper"


	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {

	cfg := config.ReadConfig("config.yml")
	logger := setupLogging(cfg)
	d := scraper.NewSpuMetricsDaemon(cfg, logger)
	// Currently config is not read correctly
	trans,_ := d.ExecuteScrape()
	// prometheus format
	reg := createMetricLines(trans)
	writeToFile(reg, cfg.Prometheus.Outfile)

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
