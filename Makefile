.PHONY: all build install uninstall clean

# Detect OS - Windows (Git Bash) vs Unix
UNAME_S := $(shell uname -s)
ifneq (,$(findstring MINGW,$(UNAME_S)))
	BINARY := ghtools.exe
else ifneq (,$(findstring MSYS,$(UNAME_S)))
	BINARY := ghtools.exe
else ifneq (,$(findstring CYGWIN,$(UNAME_S)))
	BINARY := ghtools.exe
else
	BINARY := ghtools
endif

# Install prefix: Windows → ~/bin, Unix → ~/.local/bin
ifneq (,$(findstring .exe,$(BINARY)))
	PREFIX := $(HOME)/bin
else
	PREFIX := $(HOME)/.local/bin
endif

all: build

build:
	go build -o $(BINARY) .

install: build
	@mkdir -p $(PREFIX)
	@cp $(BINARY) $(PREFIX)/$(BINARY)
	@echo "Installed ghtools to $(PREFIX)/$(BINARY)"

uninstall:
	@rm -f $(PREFIX)/ghtools $(PREFIX)/ghtools.exe
	@echo "Uninstalled ghtools from $(PREFIX)"

clean:
	rm -f ghtools ghtools.exe
