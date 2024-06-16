export PATH := $(GOPATH)/bin:$(PATH)
LDFLAGS := -s -w

all: fmt build check-size

build: gssh nssh

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

gci:
	gci write -s standard -s default -s "prefix(github.com/gofrp/tiny-frpc/)" ./

vet:
	go vet ./...

gssh:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/tiny-frpc ./cmd/frpc

nssh:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags nssh -o bin/tiny-frpc-ssh ./cmd/frpc

check-size:
	@echo "Checking file sizes..."
	@FRPC_SIZE=$$(stat -c%s "./bin/tiny-frpc"); \
	FRPC_SSH_SIZE=$$(stat -c%s "./bin/tiny-frpc-ssh"); \
	if [ $$FRPC_SSH_SIZE -gt $$FRPC_SIZE ]; then \
		echo "Error: tiny-frpc-ssh ($$FRPC_SSH_SIZE bytes) is larger than tiny-frpc ($$FRPC_SIZE bytes)"; \
		exit 1; \
	else \
		echo "File size check passed: tiny-frpc-ssh ($$FRPC_SSH_SIZE bytes) is not larger than tiny-frpc ($$FRPC_SIZE bytes)"; \
	fi

test: gotest

gotest:
	go test -v --cover ./...

alltest: vet gotest
	
clean:
	rm -f ./bin/tiny-frpc
	rm -f ./bin/tiny-frpc-ssh
	rm -rf ./lastversion

env:
	@go version
