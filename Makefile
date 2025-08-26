# Variables
GO = go
OUTPUT_DIR = build
BINARY_NAME = goFinder
LDFLAGS_LINUX = 
LDFLAGS_WINDOWS = -ldflags="-H=windowsgui"

# Directorios
CMD_DIR = ./cmd

# Objetivos
.PHONY: all run build build-linux build-windows clean test help

all: build

# Ejecutar el programa
run:
	$(GO) run $(CMD_DIR)

# Compilar para Linux
build-linux:
	mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS_LINUX) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Compilar para macOS
build-darwin:
	mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS_LINUX) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Compilar para Windows
build-windows:
	mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS_WINDOWS) -o $(OUTPUT_DIR)/$(BINARY_NAME).exe $(CMD_DIR)

# Compilar para todas las plataformas
build: build-linux build-darwin build-windows

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
	@echo "  make build        		- Compila para Linux, macOS y Windows"
	@echo "  make build-linux  		- Compila para Linux"
	@echo "  make build-windows 	- Compila para Windows"
	@echo "  make build-darwin   	- Compila para macOS"
	@echo "  make clean        		- Limpia artefactos"
	@echo "  make test         		- Ejecuta pruebas"