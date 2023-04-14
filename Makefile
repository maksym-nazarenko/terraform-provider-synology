
client-test:
	go test -v -race ./synology-go/...

test: client-test
