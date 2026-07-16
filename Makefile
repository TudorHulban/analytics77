info_color := \033[0;32m
no_color := \033[0m

repo_linting := .linting
golangci_lint_ver := 2.12.2

.DEFAULT_GOAL := lint

${repo_linting}:
	mkdir -p ${repo_linting}

install-golangci-lint: ${repo_linting}
	@echo -e "$(info_color)==> $@ $(no_color)"
	@if [ -f "${repo_linting}/golangci-lint" ] && "${repo_linting}/golangci-lint" --version >/dev/null 2>&1; then \
		echo "golangci-lint is already installed"; \
	else \
		echo "golangci-lint not found or not executable, installing v${golangci_lint_ver}..."; \
		rm -f "${repo_linting}/golangci-lint"; \
		curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ${repo_linting} v${golangci_lint_ver} || { echo "Failed to install golangci-lint"; exit 1; }; \
		echo "Installed golangci-lint, version: $$("${repo_linting}/golangci-lint" --version 2>&1)" || { echo "Failed to verify golangci-lint version"; exit 1; }; \
	fi

deps: install-golangci-lint
	@echo -e "$(info_color)==> $@ $(no_color)"


lint: deps
	@${repo_linting}/golangci-lint config path --config ${PWD}/.golangci.yaml
	@${repo_linting}/golangci-lint config verify --config ${PWD}/.golangci.yaml
	@${repo_linting}/golangci-lint run --config ${PWD}/.golangci.yaml

test:
	@echo -e "$(info_color)==> $@ $(no_color)"
	@go list ./... | while read pkg; do \
		echo "==> Testing $$pkg"; \
		go test -race -count=1 -shuffle=on -v $$pkg || exit $$?; \
	done

coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-local: 
	@go test -failfast -count=1 -shuffle=on -cpu=4,8 ./... -json -cover -race | tparse -smallscreen

# Run tests forcing the legacy RDTSCP implementation
test-legacy: 
	@go test -tags legacy_cpu -failfast -count=1 -shuffle=on -cpu=4,8 ./... -json -cover -race | tparse -smallscreen

install-fieldalignment: ${repo_linting}
	@echo -e "$(info_color)==> $@ $(no_color)"
	@if [ -f "${repo_linting}/fieldalignment" ]; then \
		echo "fieldalignment is already installed"; \
	else \
		echo "fieldalignment not found, installing..."; \
		rm -f "${repo_linting}/fieldalignment"; \
		GOBIN=$$(pwd)/${repo_linting} go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest || { echo "Failed to install fieldalignment"; exit 1; }; \
		echo "Installed fieldalignment successfully"; \
	fi

alignment-fix:
	@${repo_linting}/fieldalignment -fix ./...



# for when golangci lint does not work
# 	@fieldalignment ./...
# 	@gocritic check ./...
# 	@nilerr ./...
# 	@staticcheck ./...
# 	@unconvert -v ./...
# 	@ineffassign ./...
# 	@gocyclo -top 10 app
# 	@gocyclo -avg .
# 	@whitespace ./...