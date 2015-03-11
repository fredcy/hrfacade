VERSION = $(shell git describe --tags)
DATE = $(shell date -u +'%Y%m%dT%H%M')

build:
	go build -ldflags "-X main.version $(VERSION) -X main.build $(DATE)"


