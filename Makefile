BASE=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	cd $(BASE) && GOOS=linux go build -o target/main -ldflags="-s -w"

clean:
	cd $(BASE) && rm -rf target/
