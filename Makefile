test :
	./hack/test.sh

build :
	- cd cmd/apexd
	- go build

clean :
	- rm *.coverprofile

.PHONY:  clean
