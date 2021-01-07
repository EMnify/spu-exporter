package collector

import (
	"context"
	"strings"
	"time"

	"github.com/EMnify/spu-exporter/pkg/config"
	"github.com/EMnify/spu-exporter/pkg/prom"
	"github.com/EMnify/spu-exporter/pkg/transport"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type SpuMetricsDaemon struct {
	Cfg    *config.AppConfig
	logger log.Logger
	reg    *prometheus.Registry
}

func NewSpuMetricsDaemon(cfg *config.AppConfig, logger log.Logger) *SpuMetricsDaemon {
	return &SpuMetricsDaemon{
		Cfg:    cfg,
		logger: logger,
		reg:    prometheus.NewRegistry(),
	}
}

func (d *SpuMetricsDaemon) Run(ctx context.Context) error {
	prom.RegisterMetrics(d.reg)
	for {
		select {
		case <-ctx.Done():

			return nil
		default:
			scrapeStart := time.Now()
			trans, err := d.ExecuteScrape()
			if err != nil {
				return err
			}
			prom.CreateMetricLines(trans, d.reg)

			runtime := time.Since(scrapeStart)
			level.Debug(d.logger).Log("scrape_duration", runtime)
			time.Sleep(d.Cfg.ScrapeInterval - runtime)
		}
	}
}

func (d *SpuMetricsDaemon) ExecuteScrape() (*[]transport.Transport, error) {
	allinOne, _, err := executeScriptOnHost(d.Cfg.SSH.Host, d.Cfg.SSH.Port, d.Cfg.SSH.User, d.Cfg.SSH.Keyfile, d.Cfg.SSH.Command)

	if err != nil {
		level.Error(d.logger).Log("Failed to execute ssh: %s", err)
		return nil, err
	}

	lines := strings.Split(allinOne, "\\n")

	trans, err := parseLines(lines)
	if err != nil {
		level.Error(d.logger).Log("Error parsing transports: %s", err)
		return nil, err
	}
	return &trans, nil
}
