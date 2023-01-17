#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))/..
cd $ROOT_DIR

GO=/usr/local/go/bin/go

export POS_CLI_VERSION=1.0.1
export GIT_COMMIT_CLI=$(git rev-list -1 HEAD)
export BUILD_TIME_CLI=$(date +%s)

GOPATH=$(${GO} env | grep GOPATH | awk -F"\"" '{print $2}')
GOROOT=$(${GO} env | grep GOROOT | awk -F"\"" '{print $2}')
export PATH=${PATH}:${GOROOT}/bin:${GOPATH}/bin

mkdir -p ../kouros/api
protoc --go_out=../kouros/api --go_opt=paths=source_relative \
    --go-grpc_out=../kouros/api --go-grpc_opt=paths=source_relative \
    -I $ROOT_DIR/../../proto cli.proto

# Build CLI binary
lib/pnconnector/script/build_resource.sh

go mod vendor

# Build CLI
${GO} build -tags debug,ssloff -ldflags "-X cli/cmd.PosCliVersion=$POS_CLI_VERSION -X cli/cmd.GitCommit=$GIT_COMMIT_CLI -X cli/cmd.BuildTime=$BUILD_TIME_CLI -X cli/cmd.RootDir=$ROOT_DIR"
mv ./cli bin/poseidonos-cli

# Build CLI markdown and manpage documents
cd docs
rm markdown manpage -rf
mkdir -p markdown manpage
${GO} build -o ../bin/gen_md gen_md.go
${GO} build -o ../bin/gen_man gen_man.go
chmod +x ../bin/gen_md ../bin/gen_man
../bin/gen_md
../bin/gen_man

# Return to the cli directory
cd ..

