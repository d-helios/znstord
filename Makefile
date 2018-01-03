GLIDE_GO_EXECUTABLE ?= go
VERSION ?= $(shell git describe --tags)
GOOS ?= solaris
GOARCH ?= amd64

all: get build

get:
	glide up

build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} ${GLIDE_GO_EXECUTABLE} build -o bin/znstord -ldflags "-X main.version=${VERSION} -extldflags '-static'" znstord.go

build-tests:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} ${GLIDE_GO_EXECUTABLE} test -c -o bin/zfs_test zfs/zfs_test.go


install: build
	install -d ${DESTDIR}/opt/local/bin/
	install -m 755 ./bin/znstord ${DESTDIR}/opt/local/bin/znstord

upload:
	`ssh admin@${TEST_SRV} "mkdir -p /opt/local/bin/"
	`scp bin/znstord admin@${TEST_SRV}/opt/local/bin/znstord`

upload-and-restart: upload
	`ssh admin@{TEST_SRV} "sudo svcadm restart znstord"`


integration-test:
	@echo "integration-test is comming"

clean:
	rm -rf bin

.PHONY: all
