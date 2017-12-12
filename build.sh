#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit 990
fi

export HASH=`git rev-parse HEAD`
export DATE=`date '+%Y-%m-%d_%H:%M:%S%z'`
export PATH=$PATH:$GOPATH/bin
go get -u -f github.com/golang/dep/cmd/dep
go get -u -f github.com/aktau/github-release
dep ensure -no-vendor
mkdir -p bin
go build -ldflags "-X main.buildTime=$DATE -X main.gitHash=$HASH" -o bin/expmanager cmd/expmanager/*.go
go build -ldflags "-X main.buildTime=$DATE -X main.gitHash=$HASH" -race -tags NO_CUDA -o bin/expmanager-cpu-race cmd/expmanager/*.go
go build -ldflags "-X main.buildTime=$DATE -X main.gitHash=$HASH" -tags NO_CUDA -o bin/expmanager-cpu cmd/expmanager/*.go
go test -ldflags "-X command-line-arguments.TestRunMain=Use -X command-line-arguments.buildTime=$DATE -X command-line-arguments.gitHash=$HASH" -coverpkg="." -c -o bin/expmanager-cpu-run-coverage -tags 'NO_CUDA' cmd/expmanager/*.go
go test -ldflags "-X command-line-arguments.buildTime=$DATE -X command-line-arguments.gitHash=$HASH" -coverpkg="." -c -o bin/expmanager-cpu-test-coverage -tags 'NO_CUDA' cmd/expmanager/*.go
go test -ldflags "-X command-line-arguments.buildTime=$DATE -X command-line-arguments.gitHash=$HASH" -race -c -o bin/expmanager-cpu-test -tags 'NO_CUDA' cmd/expmanager/*.go
if ! [ -z ${TRAVIS_TAG+x} ]; then
    if ! [ -z ${GITHUB_TOKEN+x} ]; then
        github-release release --user karlmutch --repo MeshTest --tag ${TRAVIS_TAG} --pre-release && \
        github-release upload --user karlmutch --repo MeshTest  --tag ${TRAVIS_TAG} --name MeshTest --file bin/expmanager
    fi
fi
