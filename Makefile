.PHONY: dep
dep:
	@ go mod tidy && go mod verify

.PHONY: lint
lint:
	@ golangci-lint run --fix && echo OK

.PHONY: build
build:
	@ go build -o ./bin/names .

.PHONY: run
run:
	@ ./bin/names
