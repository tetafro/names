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

.PHONY: docker
docker:
	@ docker build -t ghcr.io/tetafro/names .

.PHONY: deploy
deploy:
	@ ansible-playbook \
	--private-key ~/.ssh/id_ed25519 \
	--inventory '${SSH_SERVER},' \
	--user ${SSH_USER} \
	./playbook.yml
