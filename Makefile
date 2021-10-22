# We assume to use amd64 arch to crosscompile

#some make env
ARMGCC=arm-linux-gnueabihf-gcc
ARMGPP=arm-linux-gnueabihf-g++
ARM64GCC=aarch64-linux-gnu-gcc
ARM64GPP=aarch64-linux-gnu-g++

#GOARCH env
ARMTARGET=arm
ARM64TARGET=arm64
AMD64TARGET=amd64

#C static env
LD_FLAGS="-static -static-libgcc -static-libstdc++"

#GO static env
GO_OS=linux
C_GO=1
GO_LDFLAGS="-linkmode external -extldflags '-static -static-libgcc -static-libstdc++' -s -w"

#MAKEFLAGS += --silent

RGB_LIBDIR=./third_party/rpi-rgb-led-matrix/lib
RGB_LIBRARY_NAME=librgbmatrix
RGB_LIBRARY=$(RGB_LIBDIR)/$(RGB_LIBRARY_NAME).a
RGB_LIBSO=$(RGB_LIBDIR)/$(RGB_LIBRARY_NAME).so.1

INCLUDE_DIR=./third_party/rpi-rgb-led-matrix/include
INCLUDE_LMH=$(INCLUDE_DIR)/led-matrix-c.h

GO_MATRIX_DIR=./pkg/matrix
GO_INCLUDE_DIR=$(GO_MATRIX_DIR)/include
GO_CLIB_ARM=$(GO_MATRIX_DIR)/lib/arm
GO_CLIB_ARM64=$(GO_MATRIX_DIR)/lib/arm64
GO_CLIB_AMD64=$(GO_MATRIX_DIR)/lib/amd64

GO_EXAMPLES_DIR=./examples

BIN_DIR=./bin

.PHONY: all
all: c-compile go-examples

.PHONY: clean

c-compile:
	cp $(INCLUDE_LMH) $(GO_INCLUDE_DIR)

	LDFLAGS=$(LD_FLAGS) $(MAKE) -C $(RGB_LIBDIR)
	cp $(RGB_LIBRARY) $(GO_CLIB_AMD64)
	cp $(RGB_LIBSO) $(GO_CLIB_AMD64)
	$(MAKE) -C $(RGB_LIBDIR) clean

	LDFLAGS=$(LD_FLAGS) CC=$(ARMGCC) CXX=$(ARMGPP) $(MAKE) -C $(RGB_LIBDIR)
	cp $(RGB_LIBRARY) $(GO_CLIB_ARM)
	cp $(RGB_LIBSO) $(GO_CLIB_ARM)
	$(MAKE) -C $(RGB_LIBDIR) clean

	LDFLAGS=$(LD_FLAGS) CC=$(ARM64GCC) CXX=$(ARM64GPP) $(MAKE) -C $(RGB_LIBDIR)
	cp $(RGB_LIBRARY) $(GO_CLIB_ARM64)
	cp $(RGB_LIBSO) $(GO_CLIB_ARM64)
	$(MAKE) -C $(RGB_LIBDIR) clean

	git add $(GO_MATRIX_DIR)/*

go-examples:
	CC=$(ARMGCC) CXX=$(ARMGPP) GOOS=$(GO_OS) GOARCH=$(ARMTARGET) CGO_ENABLED=$(C_GO) \
	go build -o $(BIN_DIR)/basic_$(ARMTARGET) -ldflags $(GO_LDFLAGS) $(GO_EXAMPLES_DIR)/basic/main.go

clean:
	rm -f $(BIN_DIR)/*
	rm -f $(GO_INCLUDE_DIR)/*
	rm -f $(GO_CLIB_AMD64)/*
	rm -f $(GO_CLIB_ARM)/*
	rm -f $(GO_CLIB_ARM64)/*
