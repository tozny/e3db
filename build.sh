#! /bin/sh

set -e
set -x

# Download all dependencies (really should use Glide)
go get -d ./...

mkdir -p build

for os in darwin linux windows; do
  for arch in 386 amd64; do
    GOOS=$os GOARCH=$arch go build -v -o build/e3db-cli-$os-$arch
  done
done
