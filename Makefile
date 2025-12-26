GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
NAME=apartment_service
ENV=dev
DB=mysql
COMMAND=version
STEPS=1

.PHONY: init
# init env
init:
	@which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.0
	@which swag || go install github.com/swaggo/swag/cmd/swag@v1.16.4

export-path:
	@export PATH=$$PATH:$$(go env GOPATH)/bin
	@source ~/.bashrc 2>/dev/null || true
	@export PATH=$$PATH:$$(go env GOPATH)/bin

.PHONY: run
# run
run: export-path
	@go fmt ./...
	@which golangci-lint >/dev/null 2>&1 || { go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.0; export PATH=$$PATH:$$(go env GOPATH)/bin; }
	@golangci-lint run ./...
	@go run ./cmd/app -name $(NAME) -version $(VERSION) -env $(ENV)

.PHONY: dev
# run app with DEV environment
dev:
	@make run ENV=dev

.PHONY: qc
# run app with QC environment
qc:
	@make run ENV=qc

.PHONY: build
# build
build:
	@mkdir -p build_$(ENV)/config && \
	cp -r config/base.yaml build_$(ENV)/config/ && \
	cp -r config/$(ENV).yaml build_$(ENV)/config/ && \
	go mod tidy && \
	go build -ldflags "-X main.Name=$(NAME) -X main.Version=$(VERSION) -X main.Env=$(ENV)" -o build_$(ENV)/app ./cmd/app

.PHONY: build-dev
# build for DEV environment
build-dev:
	@make build ENV=dev

.PHONY: build-qc
# build for QC environment
build-qc:
	@make build ENV=qc

.PHONY: build-prod
# build for PROD environment
build-prod:
	@make build ENV=prod

.PHONY: gen-swagger
# gen-swagger
gen-swagger:
	@swag init -g ./cmd/app/main.go -o ./docs/swagger --outputTypes json

.PHONY: run-swagger
# run-swagger
run-swagger:
	@make gen-swagger
	@docker run --rm -p 8888:8080 -e SWAGGER_JSON=/mnt/swagger.json -v $(shell pwd)/docs/swagger/swagger.json:/mnt/swagger.json swaggerapi/swagger-ui

.PHONY: push-swagger
# push-swagger
push-swagger:
	@make gen-swagger
	@cp docs/swagger/swagger.json docs/swagger/swagger.json.bak
	@sed -i 's|"host": ".*"|"host": "apartment-business.invinhome.com"|' docs/swagger/swagger.json.bak
	@sed -i -z 's|"schemes":\s*\[\s*"http"\s*\],|"schemes": ["https"],|' docs/swagger/swagger.json.bak
	@curl -X POST "https://wshr.invinhome.com/swagger/api/push?service=go-apartment" -H "Content-Type: application/json" --data-binary @docs/swagger/swagger.json.bak
	@rm -rf docs/swagger/swagger.json.bak

.PHONY: sonar
# sonar scan
sonar:
	@docker run --rm \
	  --network dev_default \
	  -e SONAR_HOST_URL=https://sonarqube.testk8s.rdapartment.com \
	  -e SONAR_TOKEN=squ_46a75e7e8f0bf1a8767025de4d6b6831809e46b5 \
	  -v $(shell pwd):/usr/src \
	  sonarsource/sonar-scanner-cli \
	  -Dsonar.projectKey=go-apartment-dev \
	  -Dsonar.projectName=go-apartment-dev \
	  -Dsonar.sources=. \
	  -Dsonar.language=go \
	  -Dsonar.qualitygate.wait=true \
	  -Dsonar.exclusions=**/vendor/**,**/logs_*/**,**/documents/**,**/local/**,**/*_test.go,**/*.json

.PHONY: migrate
# migrate
# make migrate ENV=dev DB=mysql COMMAND=up STEPS=1
# make migrate ENV=dev DB=mysql COMMAND=down STEPS=1
# make migrate ENV=dev DB=mysql COMMAND=force STEPS=10
# make migrate ENV=dev DB=mysql COMMAND=version
# make migrate ENV=dev DB=es COMMAND=version
# make migrate ENV=dev DB=es COMMAND=up STEPS=1
# make migrate ENV=dev DB=es COMMAND=down STEPS=1
# make migrate ENV=dev DB=es COMMAND=force STEPS=1
migrate:
	@go run cmd/migrate/main.go -name $(NAME) -env $(ENV) -db $(DB) -command $(COMMAND) -steps $(STEPS)

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help