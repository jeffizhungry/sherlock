#!/bin/bash
# -----------------------------------------------------------------------------
#
# Resources
# https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl
#
# Notes:
# => This script currently doesn't work, use generate_cert.go instead
# - openssl is a pain to work with
# - however it might be fun to debug what the differences are between the certs 
#   produced by this script vs generate_cert.go
#
# -----------------------------------------------------------------------------

# Variables
KEYFILE=key.pem
CERTFILE=cert.pem

# -------------------------------------
# Generate Certificates
# -------------------------------------

# Generate Private Key
openssl genrsa -out $KEYFILE 2048

openssl rsa -in $KEYFILE -out $KEYFILE

# Create Certificate Signing Request
openssl req -sha256 -newkey -nodes -config ssl/openssl.cnf -key $KEYFILE -out server.csr

# Sign certificate with extensions
openssl x509 -req -sha256 -days 30 -in server.csr -signkey $KEYFILE -out server.crt -extensions v3_req -extfile ssl/openssl.cnf

# Is this necessary?
cat server.crt $KEYFILE > $CERTFILE

# Clean up intermediate files
# rm server.crt server.csr
