test :
	- ./hack/test.sh
	- rm *.coverprofile

build :
	CGO_ENABLED=0 go build -o apexd cmd/apexd/main.go

.PHONY:  build
