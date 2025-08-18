# Variables
GO = go
GCC = gcc
OUTPUT_DIR = build
BINARY_NAME = goFinder
LDFLAGS_LINUX = 
LDFLAGS_WINDOWS = -ldflags="-H=windowsgui"

# Directorios
CMD_DIR = ./cmd
HOTKEY_DIR = ./core/hotkey

# Objetivos
.PHONY: all run build build-linux build-windows clean test

all: build

# Ejecutar el programa
run:
	$(GO) run $(CMD_DIR)

# Compilar para Linux
build-linux:
	$(GO) build $(LDFLAGS_LINUX) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Compilar para Windows
build-windows:
	$(GO) build $(LDFLAGS_WINDOWS) -o $(OUTPUT_DIR)/$(BINARY_NAME).exe $(CMD_DIR)

build-darwin:
	$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin $(CMD_DIR)

# Compilar para todas las plataformas
build: build-linux build-windows

# Limpiar artefactos
clean:
	rm -rf $(OUTPUT_DIR)/*

# Ejecutar pruebas
test:
	$(GO) test -v ./...

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make run          		- Ejecuta el programa"
	@echo "  make build        		- Compila para Linux y Windows"
	@echo "  make build-linux  		- Compila para Linux"
	@echo "  make build-windows 	- Compila para Windows (incluye DLL)"
	@echo "  make build-darwin   	- Compila para macOS"
	@echo "  make clean        		- Limpia artefactos"
	@echo "  make test         		- Ejecuta pruebas"