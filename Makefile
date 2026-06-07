# Variables
GO = go
OUTPUT_DIR = build
BINARY_NAME = goFinder
GIT_VERSION := $(shell git describe --tags --always 2>/dev/null || echo dev)
GIT_DIRTY := $(shell git status --porcelain 2>/dev/null)
ifneq ($(strip $(GIT_DIRTY)),)
	VERSION ?= $(GIT_VERSION)-dirty
else
	VERSION ?= $(GIT_VERSION)
endif
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)
VERSION_PACKAGE = github.com/adelylria/GoFinder/core/version
VERSION_LDFLAGS = -X '$(VERSION_PACKAGE).Version=$(VERSION)' -X '$(VERSION_PACKAGE).Commit=$(COMMIT)' -X '$(VERSION_PACKAGE).BuildDate=$(BUILD_DATE)'
LDFLAGS_LINUX = -ldflags="-s -w $(VERSION_LDFLAGS)"
LDFLAGS_WINDOWS = -ldflags="-H=windowsgui -s -w"
WINDRES = windres
WINDOWS_ICON_RC = $(CMD_DIR)/gofinder.rc
WINDOWS_ICON_SYSO = $(CMD_DIR)/gofinder_windows.syso

# Directorios
CMD_DIR = ./cmd

# Plataforma por defecto (evita cross-compile de Linux en Windows, etc.)
ifeq ($(OS),Windows_NT)
	DEFAULT_BUILD := build-windows
else
	UNAME_S := $(shell uname -s 2>/dev/null)
	DEFAULT_BUILD := build-linux
endif

# Objetivos
.PHONY: all run build build-all build-linux build-windows clean test help

all: $(DEFAULT_BUILD)

# Ejecutar el programa
run:
	$(GO) run $(CMD_DIR)

# Compilar para Linux
build-linux:
	mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS_LINUX) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Compilar para Windows
build-windows:
	powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build-windows.ps1

$(WINDOWS_ICON_SYSO): $(WINDOWS_ICON_RC) core/resource/assets/GoFinder.ico
	$(WINDRES) -O coff -F pe-x86-64 -o $(WINDOWS_ICON_SYSO) $(WINDOWS_ICON_RC)

# Compilar solo para la plataforma actual
build: $(DEFAULT_BUILD)

# Compilar para todas las plataformas soportadas (requiere toolchains de cross-compile)
build-all: build-linux build-windows

# Limpiar artefactos
clean:
	rm -rf $(OUTPUT_DIR)/*
	rm -f $(WINDOWS_ICON_SYSO)

# Ejecutar pruebas
test:
	$(GO) test -v ./...

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make run          		- Ejecuta el programa"
	@echo "  make / make build 		- Compila para la plataforma actual"
	@echo "  make build-all    		- Compila Linux y Windows"
	@echo "  make build-linux  		- Compila para Linux"
	@echo "  make build-windows 	- Compila para Windows"
	@echo "  make clean        		- Limpia artefactos"
	@echo "  make test         		- Ejecuta pruebas"
