package main

import (
	"log"
	"strings"
)

func main() {

	cfg := readConfig("config.yml")

	// Currently config is not read correctly
	allinOne, _, err := executeScriptOnHost(cfg.Ssh.Host, cfg.Ssh.Port, cfg.Ssh.User, cfg.Ssh.Keyfile, cfg.Ssh.Command)

	if err != nil {
		log.Fatal("Failed to execute ssh.")
	}

	lines := strings.Split(allinOne, "\\n")

	trans, err := parseLines(lines)
	if err != nil {
		log.Fatal("Error parsing transports")
	}

	// prometheus format
	reg := createMetricLines(trans)
	writeToFile(reg, cfg.Prometheus.Outfile)

}
