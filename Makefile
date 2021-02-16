include .env
export $(shell sed 's/=.*//' .env)

BIN=bin/
EXE=$(BIN)blink

build: $(EXE)

$(EXE): main.go ## main program compilation with linux arm flags
	GOOS=linux GOARCH=arm go build --ldflags "-s -w \
	  -X 'main.InfluxEndpoint=${INFLUX_ENDPOINT}' \
	  -X 'main.InfluxToken=${INFLUX_TOKEN}' \
	  -X 'main.InfluxOrg=${INFLUX_ORG}' \
	  -X 'main.InfluxBucket=${INFLUX_BUCKET}'" \
	  -o $@

.PHONY: deploy
deploy: ## deploys the current binary on the pi target
	scp $(EXE) ${PI_ADDR}:.

.PHONY: clean
clean: ## deletes every binary in bin folder
	rm -rf bin/*

.PHONY: test ## Tests the module code
	go test ./...
