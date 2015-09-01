PROGVERSION = $(shell git describe --tags)
DATE = $(shell date -u +'%Y%m%dT%H%M')

build:
	go build -ldflags "-X main.version $(PROGVERSION) -X main.build $(DATE)"

PACKAGE = hrfacade
DESTDIR = /var/tmp/$(PACKAGE)

install:
	install -D service $(DESTDIR)/usr/local/bin/$(PACKAGE)
	install -D --mode=640 $(PACKAGE).conf $(DESTDIR)/etc/init/$(PACKAGE).conf

VERSION = 0
ITERATION = 1

fpm:
	fpm -s dir -t rpm -n $(PACKAGE) -v $(VERSION) --iteration $(ITERATION) \
	   -C $(DESTDIR) --config-files /etc .
