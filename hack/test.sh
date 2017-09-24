#!/usr/bin/env bash
#
# Run all CNI tests
#   ./test
#   ./test -v
#
# Run tests for one package
#   PKG=./libcni ./test
#
set -e

export PATH=$GOROOT/bin:$PATH
export GO15VENDOREXPERIMENT=1

function runtest {
    bash -c "go test -covermode set $@"
}

TESTABLE="$(go list ./... | grep -v vendor | xargs echo)"
FORMATTABLE="$TESTABLE"
TEST=$TESTABLE
FMT=$FORMATTABLE

echo "Running tests with coverage profile generation..."
i=0
for t in ${TEST}; do
    runtest "-coverprofile ${i}.coverprofile ${t}"
    i=$((i+1))
done
gover
goveralls -service=travis-ci -coverprofile=gover.coverprofile -repotoken=$COVERALLS_TOKEN

echo "Checking gofmt..."
fmtRes=$(go fmt $FMT)
if [ -n "${fmtRes}" ]; then
	echo -e "go fmt checking failed:\n${fmtRes}"
	exit 255
fi

echo "Checking govet..."
vetRes=$(go vet $TEST)
if [ -n "${vetRes}" ]; then
	echo -e "go vet checking failed:\n${vetRes}"
	exit 255
fi

echo "Success"
