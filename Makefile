.PHONY: all test test_v generate lint vet fmt coverage check check-fast prepare race

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
PKGSDIRS=$(shell find -L . -type f -name "*.go")

all: prepare

travis: info vet lint check test_v coverage

coverage:
	@echo "$(OK_COLOR)Generate coverage$(NO_COLOR)"
	@./scripts/cover_multi.sh

prepare: generate fmt vet lint check test race

test_v:
	@echo "$(OK_COLOR)Test packages$(NO_COLOR)"
	@go test -cover -v ./...

test:
	@echo "$(OK_COLOR)Test packages$(NO_COLOR)"
	@go test -cover ./...

lint:
	@echo "$(OK_COLOR)Run lint$(NO_COLOR)"
	@test -z "$$(golint -min_confidence 0.3 ./... | tee /dev/stderr)"

check:
	@echo "$(OK_COLOR)Run golangci-lint$(NO_COLOR)"
	@golangci-lint run --no-config --exclude-use-default=true --max-same-issues=10 --disable=gosimple --disable=golint --enable=megacheck --enable=interfacer  --enable=goconst --enable=misspell --enable=unparam --enable=goimports --disable=errcheck --disable=ineffassign  --disable=gocyclo --disable=gas

vet:
	@echo "$(OK_COLOR)Run vet$(NO_COLOR)"
	@go vet ./...

race:
	@echo "$(OK_COLOR)Test for races$(NO_COLOR)"
	@go test -race .

fmt:
	@echo "$(OK_COLOR)Formatting$(NO_COLOR)"
	@echo $(PKGSDIRS) | xargs -I '{p}' -n1 goimports -w {p}

info:
	depscheck -totalonly -tests .
	golocc --no-vendor ./...

generate:
	@echo "$(OK_COLOR)Go generate$(NO_COLOR)"
	@go generate

tools:
	@echo "$(OK_COLOR)Install tools$(NO_COLOR)"
	GO111MODULE=off go get -u github.com/warmans/golocc
	GO111MODULE=off go get -u github.com/divan/depscheck
	# golint.go:21:14: l.LintPkg undefined (type *"github.com/golangci/lint-1".Linter has no field or method LintPkg)
	# gomnd.go:13:6: cannot use magic_numbers.Analyzer (type *"github.com/tommy-muehle/go-mnd/vendor/golang.org/x/tools/go/analysis".Analyzer)
	# 	as type *"golang.org/x/tools/go/analysis".Analyzer in slice literal
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.26.0
