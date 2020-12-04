package collector

import (
	"strings"

	"github.com/EMnify/spu-exporter/pkg/config"
	"github.com/EMnify/spu-exporter/pkg/transport"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type SpuMetricsDaemon struct {
	Cfg    *config.AppConfig
	logger log.Logger
}

func NewSpuMetricsDaemon(cfg *config.AppConfig, logger log.Logger) *SpuMetricsDaemon {
	return &SpuMetricsDaemon{
		Cfg:    cfg,
		logger: logger,
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
