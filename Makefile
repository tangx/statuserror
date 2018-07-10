test:
	vgo test -race `go list ./... | grep -v examples`
cover:
	vgo test -race -coverprofile=coverage.txt -covermode=atomic `go list ./... | grep -v examples`