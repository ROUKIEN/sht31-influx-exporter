include .env
export $(shell sed 's/=.*//' .env)

BIN=bin/
EXE=$(BIN)blink

PI_ADDR=pi@192.168.1.26

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

.PHONY: blink
blink: ## runs the process remotely
	ssh ${PI_ADDR} ./blink

.PHONY: clean
clean: ## deletes every binary in bin folder
	rm -rf bin/*

.PHONY: bdb
bdb: build deploy blink
