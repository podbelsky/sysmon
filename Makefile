ifneq ("$(wildcard .docker/.env)", "")
	include .docker/.env
endif
ifneq ("$(wildcard .docker/.env.override)", "")
	include .docker/.env.override
endif
ifneq ("$(wildcard .env)", "")
	include .env
endif
ifneq ("$(wildcard .env.override)", "")
	include .env.override
endif
export

PACKAGE = $(shell go list -m)
VERSION ?= $(shell git describe --exact-match --tags 2> /dev/null || head -1 CHANGELOG.md 2> /dev/null | cut -d ' ' -f 2)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")
COMMIT ?= $(shell git rev-parse HEAD)
LDFLAGS = -ldflags "-w -X ${PACKAGE}/internal/version.Version=${VERSION} -X ${PACKAGE}/internal/version.BuildDate=${BUILD_DATE} -X ${PACKAGE}/internal/version.Commit=${COMMIT}"
TAGS =
UTILS_COMMAND = docker build -q -f .docker/utils/Dockerfile .docker/utils | xargs -I % docker run --rm -v .:/src %

.PHONY: *

build: ## build a binary
	go build -tags '${TAGS}' ${LDFLAGS} -o bin/app

# Запуск/остановка локального окружения
up down stop:
	make -C .docker/development $@
bash-% logs-% restart-%:
	make -C .docker/development $@

run: ## Start docker development environment and run app
	docker exec -it sysmon-app /bin/bash -c "go run main.go grpc"

run-client:
	docker exec -it sysmon-app /bin/bash -c "go run main.go client"

# Запуск всех тестов
test:
	go test -tags mock,integration -race -cover ./...

test-unit:
	go test -tags mock -race -count 10 ./...

# Запуск всех тестов с выключенным кешированием результата
test-no-cache:
	go test -tags mock,integration -race -cover -count=1 ./...

# Запуск всех линетров
lint:
	${UTILS_COMMAND} golangci-lint run ${args}

lint-fix:
	make lint args=--fix

# Генерация grpc сервера и клиента на основе proto-файла
# Требует предустановленного buf (https://buf.build/docs/installation)
gen-grpc:
	${UTILS_COMMAND} buf generate -v --template api/grpc/buf.gen.yaml api/grpc

# Валидация proto спецификаций
# Требует предустановленного buf (https://buf.build/docs/installation)
lint-grpc:
	${UTILS_COMMAND} buf lint --config api/grpc/buf.yaml
