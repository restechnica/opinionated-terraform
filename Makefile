# make sure targets do not conflict with file and folder names
.PHONY: build clean test

# build the project
build:
	@go build -o ./bin/otf

# run quality assessment checks
check:
	@echo "Running gofmt ..."
	@find . -name '*.go' -not -path './.go*' | xargs gofmt -s -d -l
	@echo "Ok!"

	@echo "Running go vet ..."
	@go vet ./...
	@echo "Ok!"

	@echo "Running goimports ..."
	@find . -name '*.go' -not -path './.go*' | xargs go tool goimports -l
	@echo "Ok!"

# clean
clean:
	rm -rf bin out

# format
format:
	go fmt ./...
	go tool goimports -w .

# run the binary
run:
	./bin/otf

# run tests
test:
	mkdir -p ./out
	go test ./... -cover -v -coverprofile ./out/coverage.txt
	go tool uncover ./out/coverage.txt
	go tool cover -func=./out/coverage.txt
