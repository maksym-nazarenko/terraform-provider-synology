
test-client:
	go test -v -race ./synology-go/...

test: test-client
