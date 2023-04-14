# Based on https://gist.github.com/trosendal/d4646812a43920bfe94e

DEPLOC = https://github.com/spacemeshos/ed25519_bip32/releases/download
DEPTAG = 0.1.0
DEPLIB = libed25519_bip32

ifeq ($(OS),Windows_NT)
    MACHINE = WIN32
        DEPFN = $(DEPLIB).dll
#    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
#        MACHINE += AMD64
#    endif
#    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
#        MACHINE += IA32
#    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        MACHINE = LINUX
        DEPFN = $(DEPLIB).so
    endif
    ifeq ($(UNAME_S),Darwin)
        MACHINE = OSX
        DEPFN = $(DEPLIB).dylib
    endif
#    UNAME_P := $(shell uname -p)
#    ifeq ($(UNAME_P),x86_64)
#        MACHINE += AMD64
#    endif
#    ifneq ($(filter %86,$(UNAME_P)),)
#        MACHINE += IA32
#    endif
#    ifneq ($(filter arm%,$(UNAME_P)),)
#        MACHINE += ARM
#    endif
endif

# Download the platform-specific dynamic library we rely on
.PHONY: deps
deps:
	# silent, show errors, fail fast, follow links
	curl -sSfL $(DEPLOC)/v$(DEPTAG)/$(DEPFN) -o $(DEPFN)
