
TARGET_DIR      := bin
BUILD_NAME      := ifconfig.wiki
LIB_ROOT		:= $(shell pwd)/lib/clang
BUILD_VERSION   := $(shell date "+%Y%m%d.%H%M%S")
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )

all: release

release:
	go build -ldflags "-s -w  \
	-X 'ifconfig.wiki/config.Version=${BUILD_VERSION}' \
	-X 'ifconfig.wiki/config.BuildTime=${BUILD_TIME}' \
	-X 'ifconfig.wiki/config.CommitID=${COMMIT_SHA1}' \
	" -o ${TARGET_DIR}/${BUILD_NAME}
