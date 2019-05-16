#!/bin/bash -e

if ( find /project -maxdepth 0 -empty | read v );
then
  echo "source code must be mounted into the /project directory"
  exit -1
fi

set -e
set -o pipefail

export HASH=`git rev-parse HEAD`
export DATE=`date '+%Y-%m-%d_%H:%M:%S%z'`
export PATH=$PATH:$GOPATH/bin
go get -u -f github.com/golang/dep/cmd/dep
go get -u -f github.com/aktau/github-release
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
go install github.com/karlmutch/duat/cmd/semver
go install github.com/golang/dep/cmd/dep
dep ensure -no-vendor
export SEMVER=`semver`
TAG_PARTS=$(echo $SEMVER | sed "s/-/\n-/g" | sed "s/\./\n\./g" | sed "s/+/\n+/g")
PATCH=""
for part in $TAG_PARTS
do
    start=`echo "$part" | cut -c1-1`
    if [ "$start" == "+" ]; then
        break
    fi
    if [ "$start" == "-" ]; then
        PATCH+=$part
    fi
done
flags='-X github.com/leaf-ai/platform-services/internal/version.BuildTime="$DATE" -X github.com/leaf-ai/platform-services/internal/version.GitHash="$HASH" -X github.com/leaf-ai/platform-services/internal/version.SemVer="$SEMVER"'
flags="$(eval echo $flags)"
[ -e internal/gen/downstream ] || mkdir -p internal/gen/downstream
#[ -e vendor/github.com/leaf-ai/platform-services/internal ] || mkdir -p vendor/github.com/leaf-ai/platform-services/internal
#[ -e vendor/github.com/leaf-ai/platform-services/internal/gen ] || ln -s `pwd`/internal/gen vendor/github.com/leaf-ai/platform-services/internal/gen
protoc -Icmd/downstream -I/usr/include/google --plugin=$GOPATH/bin/protoc-gen-go --go_out=plugins=grpc:./internal/gen/downstream cmd/downstream/downstream.proto
if [ "$1" == "gen" ]; then
    exit 0
fi
mkdir -p cmd/downstream/bin
CGO_ENABLED=0 go build -ldflags "$flags" -o cmd/downstream/bin/downstream cmd/downstream/*.go
go build -ldflags "$flags" -race -o cmd/downstream/bin/downstream-race cmd/downstream/*.go
CGO_ENABLED=0 go test -ldflags "$flags" -coverpkg="." -c -o cmd/downstream/bin/downstream-run-coverage cmd/downstream/*.go
CGO_ENABLED=0 go test -ldflags "$flags" -coverpkg="." -c -o bin/downstream-test-coverage cmd/downstream/*.go
go test -ldflags "$flags" -race -c -o cmd/downstream/bin/downstream-test cmd/downstream/*.go
if [ -z "$PATCH" ]; then
    if ! [ -z "${SEMVER}" ]; then
        if ! [ -z "${GITHUB_TOKEN}" ]; then
            github-release release --user leaf-ai --repo platform-services --tag ${SEMVER} --pre-release && \
            github-release upload --user leaf-ai --repo platform-services  --tag ${SEMVER} --name platform-services --file cmd/downstream/bin/downstream
        fi
    fi
fi
