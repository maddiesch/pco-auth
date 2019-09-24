.PHONY: test
test:
	go test -v ./...

example/.config:
	cp example/.config.example example/.config
	$(error Please update example/.config with real values)

.PHONY: run
run: example/.config
	cd example && $(MAKE) run
