#!/usr/bin/env bash
set -e

# ===== Config (edit if needed) =====

CERT_DIR="."
DAYS_CA=3650
DAYS_CERT=825

# ===== Create directory =====

mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

echo "===> Generating Root CA..."
openssl genrsa -out rootCA.key 4096

openssl req -x509 -new -nodes \
-key rootCA.key \
-sha256 -days $DAYS_CA \
-out rootCA.crt \
-subj "/CN=Local Dev Root CA"


echo "===> Generating server key..."
openssl genrsa -out client.key 2048

echo "===> Generating CSR..."
openssl req -new \
-key client.key \
-out client.csr \
-config san.cnf

echo "===> Signing certificate with Root CA..."
openssl x509 -req \
-in client.csr \
-CA rootCA.crt \
-CAkey rootCA.key \
-CAcreateserial \
-out client.crt \
-days $DAYS_CERT \
-sha256 \
-extensions v3_req \
-extfile san.cnf

echo "===> Generating DH parameters (this may take a while)..."
openssl dhparam -out dhparam.pem 2048

echo "===> Verifying certificate SAN..."
openssl x509 -in client.crt -text | grep -A1 "Subject Alternative Name"

echo ""
echo "✅ Done! Generated files in $CERT_DIR:"
echo " - rootCA.key"
echo " - rootCA.crt"
echo " - client.key"
echo " - client.crt"
echo " - dhparam.pem"
