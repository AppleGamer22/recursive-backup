# done -> CompletionStatus
# json -> target
# target usb, source smb

VERSION:=1.0.0
SOURCES:=./internal/*
BIN:=./bin
BASE:=rb_$(VERSION)
LINUX=$(BASE)_linux_amd64
DARWIN=$(BASE)_darwin_amd64
WINDOWS=$(BASE)_windows_amd64.exe

all: build test

build: $(LINUX) $(DARWIN) $(WINDOWS)

test:
		@echo
		@echo "Testing $(VERSION)"
		go test -v -cover $(SOURCES)

linux: $(LINUX)

darwin: $(DARWIN)

windows: $(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(BIN)/$(LINUX) $(EXSOURCES)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BIN)/$(DARWIN) $(EXSOURCES)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build  -v -o $(BIN)/$(WINDOWS) $(EXSOURCES)

.PHONY: build test benchmark build-example test-integ