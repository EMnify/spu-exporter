# spu-exporter
Prometheus exporter for spu nodes, utilizing ssh interface of spu application

## Description

This application consists of three parts:
 - connecting over ssh, executing a command a returning the stdout as string.
 - parsing this output expecting specific keywords like returned by [spu application](https://www.applicata.bg/jnspu.html)
 - creating prometheus metrics and writing to a file to be collected by node exporter
 
Currently it is doing it once on call. This behaviour might change in future versions to execute it in a specified interval.

## Configuration

Configurable are:
 - the command to be executed on ssh session
 - ssh connection parameters
 - prometheus metrics host and port
 - spu scrape interval

## Execution

The path to the config file can be given as first parameter, otherwise it uses /opt/spu/exporter-config

## Project structure

When there is a project structure created, it shall be documented here :)