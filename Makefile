test :
	- ./hack/test.sh

build :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apexd cmd/apexd/main.go

clean :
	rm *.coverprofile

.PHONY: test build clean
