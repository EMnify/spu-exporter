package scraper

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

func (d *SpuMetricsDaemon) ExecuteScrape() (*[]transport.Transport,error){
	allinOne, _, err := executeScriptOnHost(d.Cfg.Ssh.Host, d.Cfg.Ssh.Port, d.Cfg.Ssh.User, d.Cfg.Ssh.Keyfile, d.Cfg.Ssh.Command)

	if err != nil {
		level.Error(d.logger).Log("Failed to execute ssh.")
	}

	lines := strings.Split(allinOne, "\\n")

	trans, err := parseLines(lines)
	if err != nil {
		level.Error(d.logger).Log("Error parsing transports")
		return nil, err
	}
	return &trans, nil
}