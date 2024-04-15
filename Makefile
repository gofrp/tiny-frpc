export PATH := $(GOPATH)/bin:$(PATH)
LDFLAGS := -s -w

all: fmt build

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
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags gssh -o bin/tiny-frpc ./cmd/frpc

nssh:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags nssh -o bin/tiny-frpc-ssh ./cmd/frpc

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