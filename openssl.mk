export define configure_csr
cat <<EOF > internal/services/pki/testdata/certs/test.needle.local.cnf
[ req ]
prompt = no
distinguished_name = dn
req_extensions = req_ext

[ dn ]
CN = test.needle.local
O = Pixelfactory
C = FR

[ req_ext ]
subjectAltName = DNS: localhost, DNS: test.needle.local, IP: 127.0.0.1
EOF
endef

testdata/certs:
	@mkdir -p internal/services/pki/testdata/certs

data/certs:
	@mkdir -p data/certs

data/certs/root-ca.key: data/certs
	@openssl genrsa -out data/certs/root-ca.key

data/certs/root-ca.crt: data/certs/root-ca.key
	@openssl req -x509 -new -nodes -sha256 -days 1825 \
		-subj '/CN=identity.needle.local/O=Pixelfactory/C=FR' \
		-key data/certs/root-ca.key \
		-out data/certs/root-ca.crt

ca: data/certs/root-ca.crt data/certs/root-ca.key
.PHONY: ca

internal/services/pki/testdata/certs/root-ca.key: testdata/certs
	@openssl genrsa -out internal/services/pki/testdata/certs/root-ca.key

internal/services/pki/testdata/certs/root-ca.crt: internal/services/pki/testdata/certs/root-ca.key 
	@openssl req -x509 -new -nodes -sha256 -days 1825 \
		-subj '/CN=identity-test.needle.local/O=Pixelfactory/C=FR' \
		-key internal/services/pki/testdata/certs/root-ca.key \
		-out internal/services/pki/testdata/certs/root-ca.crt

ca-test: internal/services/pki/testdata/certs/root-ca.crt internal/services/pki/testdata/certs/root-ca.key
.PHONY: ca-test

internal/services/pki/testdata/certs/test.needle.local.key:
	@openssl genrsa -out internal/services/pki/testdata/certs/test.needle.local.key

internal/services/pki/testdata/certs/test.needle.local.cnf:
	@bash -c "eval \"$$configure_csr\""

internal/services/pki/testdata/certs/test.needle.local.csr: internal/services/pki/testdata/certs/test.needle.local.cnf internal/services/pki/testdata/certs/test.needle.local.key
	@openssl req -new \
		-config internal/services/pki/testdata/certs/test.needle.local.cnf \
		-key internal/services/pki/testdata/certs/test.needle.local.key \
		-out internal/services/pki/testdata/certs/test.needle.local.csr

internal/services/pki/testdata/certs/test.needle.local.crt: internal/services/pki/testdata/certs/test.needle.local.csr
	@openssl x509 -req -days 825 -sha256 \
		-CA internal/services/pki/testdata/certs/root-ca.crt \
		-CAkey internal/services/pki/testdata/certs/root-ca.key \
		-CAcreateserial \
		-in internal/services/pki/testdata/certs/test.needle.local.csr \
		-out internal/services/pki/testdata/certs/test.needle.local.crt \

cert-test: internal/services/pki/testdata/certs/test.needle.local.crt internal/services/pki/testdata/certs/test.needle.local.key
.PHONY: cert-test
