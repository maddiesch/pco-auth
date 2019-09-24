.PHONY: test
test:
	go test -v ./...

example/.config:
	cp example/.config.example example/.config
	$(error Please update example/.config with real values)

.PHONY: run
run: example/.config
	cd example && $(MAKE) run

.PHONY: release
release:
	ifeq ($(git diff --quiet || echo 'dirty'),dirty)
		$(error Working directory is dirty)
	endif

	# mkdir releases
	# GOOS=darwin GOARCH=amd64 cd example && $(MAKE) build
	# cp example/build/example releases/pco-auth-darwin-amd64
	# GOOS=linux GOARCH=amd64 cd example && $(MAKE) build
	# cp example/build/example releases/pco-auth-linux-amd64
