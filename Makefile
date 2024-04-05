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
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags frpc-gssh -o bin/frpc-gssh ./cmd/go_ssh

nssh:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags frpc-nssh -o bin/frpc-nssh ./cmd/native_ssh

test: gotest

gotest:
	go test -v --cover ./...

alltest: vet gotest
	
clean:
	rm -f ./bin/gssh
	rm -f ./bin/nssh
	rm -rf ./lastversion

env:
	@go version