.PHONY: all build install uninstall clean

# Detect Windows: OS=Windows_NT is set on all Windows environments
ifeq ($(OS),Windows_NT)
BINARY := ghtools.exe
PREFIX := $(USERPROFILE)/.local/bin
else
BINARY := ghtools
PREFIX := $(HOME)/.local/bin
endif

all: build

build:
	go build -o $(BINARY) .

install: build
	mkdir -p "$(PREFIX)"
	cp $(BINARY) "$(PREFIX)/$(BINARY)"
	@echo "Installed ghtools to $(PREFIX)/$(BINARY)"

uninstall:
	rm -f "$(PREFIX)/ghtools" "$(PREFIX)/ghtools.exe"
	@echo "Uninstalled ghtools from $(PREFIX)"

clean:
	rm -f ghtools ghtools.exe
