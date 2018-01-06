#!/bin/bash -x

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit 990
fi

export HASH=`git rev-parse --short HEAD`
export DATE=`date '+%Y-%m-%d_%H:%M:%S%z'`
export PATH=$PATH:$GOPATH/bin
go get -u -f github.com/golang/dep/cmd/dep
go get -u -f github.com/aktau/github-release
dep ensure -no-vendor
mkdir -p cmd/restpoc/bin
go build -ldflags "-X github.com/SentientTechnologies/platform-services/version.BuildTime=$DATE -X github.com/SentientTechnologies/platform-services/version.GitHash=$HASH" -o cmd/restpoc/bin/restpoc cmd/restpoc/*.go
go build -ldflags "-X github.com/SentientTechnologies/platform-services/version.BuildTime=$DATE -X github.com/SentientTechnologies/platform-services/version.GitHash=$HASH" -race -o cmd/restpoc/bin/restpoc-race cmd/restpoc/*.go
go test -ldflags "-X github.com/SentientTechnologies/platform-services/version.TestRunMain=Use -X github.com/SentientTechnologies/platform-services/version.BuildTime=$DATE -X github.com/SentientTechnologies/platform-services/version.GitHash=$HASH" -coverpkg="." -c -o cmd/restpoc/bin/restpoc-run-coverage cmd/restpoc/*.go
go test -ldflags "-X github.com/SentientTechnologies/platform-services/version.BuildTime=$DATE -X github.com/SentientTechnologies/platform-services/version.GitHash=$HASH" -coverpkg="." -c -o cmd/restpoc/bin/restpoc-test-coverage cmd/restpoc/*.go
go test -ldflags "-X github.com/SentientTechnologies/platform-services/version.BuildTime=$DATE -X github.com/SentientTechnologies/platform-services/version.GitHash=$HASH" -race -c -o cmd/restpoc/bin/restpoc-test cmd/restpoc/*.go
if ! [ -z ${TRAVIS_TAG+x} ]; then
    if ! [ -z ${GITHUB_TOKEN+x} ]; then
        github-release release --user SentientTechnologies --repo platform-services --tag ${TRAVIS_TAG} --pre-release && \
        github-release upload --user SentientTechnologies --repo platform-services  --tag ${TRAVIS_TAG} --name platform-services --file cmd/restpoc/bin/restpoc
    fi
fi