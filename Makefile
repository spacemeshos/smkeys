# Based on https://gist.github.com/trosendal/d4646812a43920bfe94e

DEPURL = https://github.com/spacemeshos/spacemesh-sdk/releases/download/
DEPTAG = 0.0.1
DEPLIB = spacemesh-sdk
DEPDIR = deps
REALDEPDIR = $(shell realpath $(DEPDIR))

ifeq ($(OS),Windows_NT)
#	MACHINE = WIN32
	DEPFN = windows-amd64
#	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
#		MACHINE += AMD64
#	endif
#	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
#		MACHINE += IA32
#	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		MACHINE = linux
		CGO_LDFLAGS := -Wl,-rpath,$$ORIGIN
	endif
	ifeq ($(UNAME_S),Darwin)
		MACHINE = macos
		CGO_LDFLAGS := -Wl,-rpath,@loader_path
	endif

	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		PLATFORM = $(MACHINE)-amd64
	endif
#	ifneq ($(filter %86,$(UNAME_P)),)
#		MACHINE += IA32
#	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		PLATFORM = $(MACHINE)-arm64
	endif
	ifeq ($(UNAME_P), aarch64)
		PLATFORM = $(MACHINE)-arm64
	endif
endif
FN = $(DEPLIB)_$(PLATFORM).tar.gz

# Download the platform-specific dynamic library we rely on
.PHONY: deps
deps:
	@mkdir -p $(DEPDIR)
	@# silent, show errors, fail fast, follow links
	curl -sSfL $(DEPURL)/v$(DEPTAG)/$(FN) -o deps/$(FN)
	cd $(DEPDIR) && tar -xzf $(FN) --exclude=LICENSE

.PHONY: test
test:
	CGO_CFLAGS="-I$(REALDEPDIR)" \
	CGO_LDFLAGS="-L$(REALDEPDIR) $(CGO_LDFLAGS)" \
	LD_LIBRARY_PATH=$(REALDEPDIR) \
	go test -v -count 1 ./...

.PHONY: go-env
go-env:
	go env -w CGO_CFLAGS="-I$(REALDEPDIR)"
	go env -w CGO_LDFLAGS="-L$(REALDEPDIR) $(CGO_LDFLAGS)"
