# meta-suite v2 kernel. `make check` is the gate to run before every push (D-052); remote CI is advisory.
#   make build   5 static binaries into ../bin/ (CGO_ENABLED=0, -trimpath -s -w), sizes printed
#   make test    go test (unit, end-to-end, sheaf parity)
#   make check   gofmt + go vet + go test, then `meta check --skip-go` on the repo (spec + hygiene checks)
GO       ?= go
BIN      := ../bin
VERSION  ?= 2.0.0-alpha.1
BUILD    := $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo unset)
LDFLAGS  := -s -w -X main.version=$(VERSION) -X main.buildHash=$(BUILD)
PLATFORMS := linux/amd64 linux/arm64 darwin/arm64 darwin/amd64 windows/amd64
export CGO_ENABLED=0

.PHONY: build test vet fmt check clean host

build:
	@mkdir -p $(BIN)
	@for t in $(PLATFORMS); do \
	  os=$${t%/*}; arch=$${t#*/}; ext=""; [ "$$os" = windows ] && ext=".exe"; \
	  out="$(BIN)/meta-$$os-$$arch$$ext"; \
	  GOOS=$$os GOARCH=$$arch $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o "$$out" ./cmd/meta || exit 1; \
	  printf '%-24s %9d bytes\n' "$$os/$$arch" "$$(wc -c < "$$out")"; \
	done

# The binary for this machine only (what `check` uses).
host:
	@mkdir -p $(BIN)
	@os=$$($(GO) env GOOS); arch=$$($(GO) env GOARCH); ext=""; [ "$$os" = windows ] && ext=".exe"; \
	  $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o "$(BIN)/meta-$$os-$$arch$$ext" ./cmd/meta

fmt:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet:
	$(GO) vet ./...

test:
	$(GO) test -count=1 ./...

check: fmt vet test host
	@sh ../bin/meta check --skip-go --repo ../.. > /dev/null 2>&1 || { sh ../bin/meta check --skip-go --repo ../..; exit 1; }
	@echo "check: OK (gofmt, go vet, go test incl. sheaf parity, meta check: SKILL.md words, gitignore, atom count, D-### citations)"

clean:
	rm -f $(BIN)/meta-*
