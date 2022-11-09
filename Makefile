.PHONY: unit-tests
unit-tests:
	go test -race -short -v --count 1 ./...

