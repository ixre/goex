#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o gof-gen gof-gen.go
CGO_ENABLED=0 GOOS=darwin ARCH=amd64 go build -o mac-gof-gen gof-gen.go
CGO_ENABLED=0 GOOS=windows ARCH=amd64 go build -o gof-gen.exe gof-gen.go

tar cvzf generator-build-bin.tar.gz gof-gen.sh mac-gof-gen gof-gen\
 gof-gen.exe gen.conf templates README.md

rm -rf gof-gen mac-gof-gen gof-gen.exe