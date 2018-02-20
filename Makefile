NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

IGNORED_PACKAGES := /vendor/

deps:
	@echo "$(OK_COLOR)==> Installing golint$(NO_COLOR)"
	@go get -u github.com/golang/lint/golint
	@echo "$(OK_COLOR)==> Installing go-swagger$(NO_COLOR)"
	@type swagger >/dev/null 2>&1 || brew tap go-swagger/go-swagger && brew install go-swagger
	@echo "$(OK_COLOR)==> Installing multi-file-swagger$(NO_COLOR)"
	@type multi-file-swagger >/dev/null 2>&1 || npm install -g multi-file-swagger

test:
	@/bin/sh -c "./test.sh $(allpackages)"

run:
	go run service/cmd/service1-server/main.go --port 8001

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/Sef1995/zipkin-swagger-poc/service/cmd/service1-server

build-docker:
	docker build -t service1 .

run-docker:
	docker run -p 8001:8001 service1

build-run-docker:
	make build-docker && make run-docker

generate:
	multi-file-swagger swagger.yml > /tmp/swagger.json
	cd service/ && find . ! -name \configure_*.go -type f -exec rm -f {} +
	cd service/ && find . -type d -empty -delete
	cd service/ && swagger generate server -A service1 -f /tmp/swagger.json
	cd service/ && swagger generate client -A service1 -f /tmp/swagger.json
	git add .

_allpackages = $(shell ( go list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)