all: tests build


HOSTNAME := $(shell hostname)
ifeq ($(HOSTNAME), TerminusEst.local)
		GOPATH := /Users/bstar/workspace/goprojects/
endif


COMPILER = GOPATH=$(GOPATH) go
BUILDER = $(COMPILER) build
TEST_COMPILER = $(COMPILER) test
DEPENDER = $(COMPILER) get

ROOT = github.com/Kerah/goaxer
APPS_PATH = $(ROOT)/apps

CORE_PACKAGES = $(ROOT)/imdg

VERSION = 0.0.1
TESTING_PACKAGES = $(CORE_PACKAGES)

build:
		echo "builded"

tests:
		$(TEST_COMPILER) $(TESTING_PACKAGES) -cover

deps:
		$(DEPENDER) gopkg.in/redis.v1