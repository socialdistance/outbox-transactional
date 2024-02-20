PORT?=9999
APP_NAME?=test-app
MOQPATH  := $(shell pwd)/bin/

clean:
	rm -f ${APP_NAME}

build: clean
	go build -o ${APP_NAME}

run: build
	PORT=${PORT} ./${APP_NAME}

runx:
	go run cmd/main.go --conf=config.yaml

runx-kafka:
	go run cmd/cron/outbox_producer/outbox_producer.go  --conf=config.yaml

kafka:
	/home/user/kafka_2.13-3.6.1/bin/kafka-console-consumer.sh --topic quickstart-events --from-beginning --bootstrap-server localhost:9092

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	go generate internal/...
	mockgen -source=internal/pkg/repository/order/repository.go \
	-destination=internal/pkg/repository/order/mocks/mock_repository.go

PHONE: lint
lint:
	@(golangci-lint run)