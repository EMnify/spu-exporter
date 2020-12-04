package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/EMnify/spu-exporter/pkg/collector"
	"github.com/EMnify/spu-exporter/pkg/config"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
)

var (
	// Version gets defined by the build system.
	Version = "0.0.0"
	// Revision gets defined by the built system
	Revision = ""
	// BuildDate defines the date this binary was built.
	BuildDate string
	// GoVersion running this binary.
	GoVersion = runtime.Version()
)

func main() {

	cfg := config.ReadConfig("configs/config.yml")
	logger := setupLogging(cfg)

	app := &cli.App{
		Name:    "FritzExporter",
		Version: fmt.Sprintf("%s (%s)", Version, Revision),
		Usage:   "FritzExporter",
	}

	app.Action = func(c *cli.Context) error {
		return execute(cfg, logger)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
func execute(cfg *config.AppConfig, logger log.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := collector.NewSpuMetricsDaemon(cfg, logger)
	var g run.Group
	{
		logger = log.With(logger, "component", "fritz_exporter")
		level.Info(logger).Log(
			"msg", "starting fritzbox exporter",
			"version", Version,
			"revision", Revision,
			"buildDate", BuildDate,
			"goVersion", GoVersion,
		)

		// Currently config is not read correctly
		g.Add(func() error {
			return d.Run(ctx)
		}, func(_ error) {
			level.Info(logger).Log("msg", "shutting down socket server")
		})
	}
	{
		logger := setupLogging(cfg)
		logger = log.With(logger, "component", "metrics")

		promauto.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "spu_exporter_build_info",
				Help: "A metric with a constant '1' value labeled by version, revision, build_date, and goversion.",
				ConstLabels: prometheus.Labels{
					"version":    Version,
					"revision":   Revision,
					"build_date": BuildDate,
					"goversion":  GoVersion,
				},
			},
			func() float64 { return 1 },
		)
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		s := http.Server{
			Addr:    cfg.Prometheus.Host + ":" + cfg.Prometheus.Port,
			Handler: m,
		}
		g.Add(func() error {
			level.Info(logger).Log("msg", "starting metrics server", "addr", cfg.Prometheus.Host+":"+cfg.Prometheus.Port)
			return s.ListenAndServe()
		}, func(_ error) {
			level.Info(logger).Log("msg", "shutting down metric server")
			if err := s.Shutdown(context.Background()); err != nil {
				level.Error(logger).Log("msg", "error shutting down metrics server", "error", err)
			}
		})
	}

	{
		sig := make(chan os.Signal)
		g.Add(func() error {
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig
			return nil
		}, func(err error) {
			cancel()
			close(sig)
		})
	}
	if err := g.Run(); err != nil {
		return err
	}
	return nil
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
