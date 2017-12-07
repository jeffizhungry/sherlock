###############################################################################
#
#	Welcome to the most awesome Makefile you will ever see 😏
#
###############################################################################

#--------------------------------------
# Variables
#--------------------------------------
APP := sherlock
PROJ := github.com/jeffizhungry/sherlock

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

#--------------------------------------
# Targets
#--------------------------------------
$(APP): $(SOURCES)
	go build -o $(APP) .

#--------------------------------------
# Basic Rules
#--------------------------------------
all: $(APP)

clean:
	rm -f $(APP)

#--------------------------------------
# Custom Rules
#--------------------------------------

# Run app in foreground
.PHONY: run
run: $(APP)
	./$(APP)

# Run in dev mode with fswatch to restart on file changes
.PHONY: dev
dev:
	fswatch

# Build ctags
.PHONY: ctags
ctags:
	gotags -tag-relative=true -R=true -sort=true -f=tags -fields=+l .

# Startup test http server
.PHONY: http
http:
	go run ./examples/basic_http.go

#--------------------------------------
# Testing Rules
#--------------------------------------

.PHONY: testall
testall: clean all test vet errcheck

.PHONY: vet
vet:
	go vet $(PROJ)

.PHONY: errcheck
errcheck:
	errcheck $(PROJ)

.PHONY: test
test:
	go test -v -timeout 30s $(PROJ)
