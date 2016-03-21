VERSION = $(shell cat .version)
GHT = $(GITHUB_TOKEN)

all: 
	@echo all binaries built
	@echo Building firehose
	go build -o firehose/firehose -ldflags "-X main.VERSION=${VERSION}" firehose/main.go 
	@echo Building firetruck
	go build -o firetruck/firetruck -ldflags "-X main.VERSION=${VERSION}" firetruck/main.go 

release: firehose firetruck
	mkdir -p dist/usr/sbin
	mv eventilator dist/usr/sbin
	mv firetruck/firetruck dist/usr/sbin
	mv firehost/firehost dist/usr/sbin
	cd dist && tar -cvzf ../slammer-${VERSION}.tar.gz usr/ && cd ..
	echo Version=${VERSION}
	ls -lh eventilator-${VERSION}.tar.gz
	@echo  ghr --username therealbill --token NOTSHOWN ${VERSION} eventilator-${VERSION}.tar.gz
	ghr  --username therealbill --token ${GHT} ${VERSION} eventilator-${VERSION}.tar.gz

