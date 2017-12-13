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
	make keypair

$(KEYFILE):
	make keypair

#--------------------------------------
# Basic Rules
#--------------------------------------
all: $(APP)

clean:
	rm -f $(APP)

distclean:
	rm -f $(APP) $(CERTFILE) $(KEYFILE) server.crt server.csr
	rm -rf data

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
	go run ./developer/basic_http.go

# Generate self-signed key pair
#
# Resources:
# https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl
.PHONY: keypair
keypair:
	openssl genrsa -out $(KEYFILE) 2048
	openssl rsa -in $(KEYFILE) -out $(KEYFILE)
	# Create Certificate Signing Request
	openssl req -sha256 -newkey -nodes -config ssl/openssl.cnf -key $(KEYFILE) -out server.csr
	# Sign certificate with extensions
	openssl x509 -req -sha256 -days 30 -in server.csr -signkey $(KEYFILE) -out server.crt -extensions v3_req -extfile ssl/openssl.cnf
	cat server.crt $(KEYFILE) > $(CERTFILE)
	# Clean up intermediate files
	# rm server.crt server.csr

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
