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
CERTFILE := cert.pem
KEYFILE := key.pem

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

#--------------------------------------
# Targets
#--------------------------------------
$(APP): $(SOURCES)
	go build -o $(APP) .

$(CERTFILE):
	make cert

$(KEYFILE):
	make cert

#--------------------------------------
# Basic Rules
#--------------------------------------
all: $(APP)

clean:
	rm -f $(APP)

distclean:
	rm -f $(APP) $(CERTFILE) $(KEYFILE)
	rm -rf data

#--------------------------------------
# Custom Rules
#--------------------------------------

# Run app in foreground
.PHONY: run
run: $(APP) $(CERTFILE) $(KEYFILE)
	./$(APP)

# Run in dev mode with fswatch to restart on file changes
.PHONY: dev
dev: $(APP) $(CERTFILE) $(KEYFILE)
	fswatch

# Build ctags
.PHONY: ctags
ctags:
	gotags -tag-relative=true -R=true -sort=true -f=tags -fields=+l .

# Startup test http server
.PHONY: http
http:
	go run ./developer/basic_http.go

# Generate self-signed cert
.PHONY: cert
cert:
	go run ssl/generate_cert.go -ca

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
