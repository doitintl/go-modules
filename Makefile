-include .env

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.DEFAULT_GOAL := default

WDIR := $(shell pwd)

test-worker:
	cd worker; \
	go test -race;
