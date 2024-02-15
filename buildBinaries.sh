#!/bin/bash
version=1.1.1

mkdir build/
rm build/*

# Windows amd64
goos=windows
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o tinja.exe
zip build/TInjA_"$version"_"$goos"_"$goarch".zip tinja.exe

# Linux amd64
goos=linux
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o tinja
tar cfvz build/TInjA_"$version"_"$goos"_"$goarch".tar.gz tinja

# Linux arm64
goos=linux
goarch=arm64
GOOS=$goos GOARCH=$goarch go build -o tinja
tar cfvz build/TInjA_"$version"_"$goos"_"$goarch".tar.gz tinja

# Darwin/MacOS amd64
goos=darwin
goarch=amd64
GOOS=$goos GOARCH=$goarch go build -o tinja
tar cfvz build/TInjA_"$version"_"$goos"_"$goarch".tar.gz tinja

# Darwin/MacOS arm64
goos=darwin
goarch=arm64
GOOS=$goos GOARCH=$goarch go build -o tinja
tar cfvz build/TInjA_"$version"_"$goos"_"$goarch".tar.gz tinja

# reset
GOOS=
GOARCH=

# remove binaries
rm tinja
rm tinja.exe

# generate checksum file
find build/ -type f  \( -iname "*.tar.gz" -or -iname "*.zip" \) -exec sha256sum {} + > build/TInjA_"$version"_checksums_sha256.txt
