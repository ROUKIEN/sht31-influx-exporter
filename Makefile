BIN=bin/
EXE=$(BIN)blink

PI_ADDR=pi@192.168.1.26

build: $(EXE)

$(EXE): main.go ## main program compilation with linux arm flags
	GOOS=linux GOARCH=arm go build -o $@ $<

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