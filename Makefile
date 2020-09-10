# Go parameters
GOCMD				=go
GOBUILD				=$(GOCMD) build
GOINSTALL			=$(GOCMD) install
GOCLEAN				=$(GOCMD) clean
GOTEST				=$(GOCMD) test
TEST_FLAGS        	?=-v
BINARY_NAME			?=api
BIN					=$$PWD/bin
GOMODULES			?=on
VCS_REF	 			?=$(shell git describe --tags --long --abbrev=8 --always HEAD)
GOFLAGS        		?=-mod=vendor -gcflags='-e' -ldflags "-X main.build=${VCS_REF}"

PG_HOST           	?= postgres-db
PG_PORT           	?= 5432
PG_PASSWORD       	?= postgres
PG_USER           	?= postgres
PG_DATABASE       	?= apimain

##
## Build
##
define build
	GOGC=off GO111MODULE=$(GOMODULES) GOBIN=$(BIN) CGO_ENABLED=0 \
	$(GOINSTALL) -v $(GOFLAGS) $(1)
endef

build:
	$(call build,./cmd/$(BINARY_NAME))

##
## Run
##
define run
	sh -c '$(MAKE) build && \
	./bin/$(1)'
endef

clean: 
	$(GOCLEAN) -modcache -testcache
	rm -rf ./bin/$(BINARY_NAME)

run:
	$(call run,$(BINARY_NAME))

migrate:
	$(call run,migrate)

##
## Database
##
db-create:
	@env PG_PASSWORD=$(PG_PASSWORD) PG_USER=$(PG_USER) PG_HOST=$(PG_HOST) ./scripts/db.sh create $(PG_DATABASE)

db-drop:
	@env PG_PASSWORD=$(PG_PASSWORD) PG_USER=$(PG_USER) PG_HOST=$(PG_HOST) ./scripts/db.sh drop $(PG_DATABASE)

db-reset: db-create migrate

##
## Tools
##
tools:
	export GO111MODULE=off && \
	go get -u github.com/goware/modvendor
	go get -u golang.org/x/tools/cmd/goimports

##
## Dependency mgmt
##
dep:
	@export GO111MODULE=on && \
		go mod tidy && \
		rm -rf ./vendor && go mod vendor && \
		modvendor -copy="**/*.c **/*.h **/*.s **/*.proto"

dep-upgrade-all:
	@GO111MODULE=on go get -u ./...
	@$(MAKE) dep