
generate:
	go generate ./...

test-client:
	go test -v ./synology-go/...

test: test-client

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
