include .config

.PHONY: build
build:
	go build -o build/example .

.PHONY: run
run: .config build
	build/example -port 8080 -client_id ${PCO_CLI_CLIENT_ID} -client_secret ${PCO_CLI_CLIENT_SECRET} -scope people,services
