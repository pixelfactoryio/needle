testdata/certs:
	@mkdir -p testdata/certs

data/certs:
	@mkdir -p data/certs

data/certs/root-ca.crt data/certs/root-ca.key: data/certs
	@step certificate create identity.needle.local \
		data/certs/root-ca.crt \
		data/certs/root-ca.key \
		--profile root-ca \
		--not-after=87600h \
		--no-password \
		--insecure \
		--kty RSA

ca: data/certs/root-ca.crt data/certs/root-ca.key
.PHONY: ca

testdata/certs/root-ca.crt testdata/certs/root-ca.key: testdata/certs
	@step certificate create identity-test.needle.local \
		testdata/certs/root-ca.crt \
		testdata/certs/root-ca.key \
		--profile root-ca \
		--not-after=87600h \
		--no-password \
		--insecure \
		--kty RSA

ca-test: testdata/certs/root-ca.crt testdata/certs/root-ca.key
.PHONY: ca-test

testdata/certs/test.needle.local.crt testdata/certs/test.needle.local.key: testdata/certs
	@step certificate create test.needle.local \
		testdata/certs/test.needle.local.crt \
		testdata/certs/test.needle.local.key \
		--ca testdata/certs/root-ca.crt \
		--ca-key testdata/certs/root-ca.key \
		--profile leaf \
		--not-after 8760h \
		--no-password \
		--insecure \
		--san 127.0.0.1 \
		--san localhost \
		--san test.needle.local

cert-test: testdata/certs/test.needle.local.crt testdata/certs/test.needle.local.key
.PHONY: cert-test
