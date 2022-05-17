# spu-exporter
Prometheus exporter for spu nodes, utilizing ssh interface of spu application

## Description

This application consists of three parts:
 - connecting over ssh, executing a command a returning the stdout as string.
 - parsing this output expecting specific keywords like returned by [spu application](https://www.applicata.bg/jnspu.html)
 - creating prometheus metrics and serve them on configurable port
 
## Configuration

Configurable are:
 - the command to be executed on ssh session
 - ssh connection parameters
 - prometheus metrics host and port
 - spu scrape interval

## Execution

The path to the config file can be given as first parameter, otherwise it uses /opt/spu/exporter-config.yaml

