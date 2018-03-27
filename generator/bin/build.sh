#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o gof-gen_linux gof-gen.go
CGO_ENABLED=0 GOOS=darwin ARCH=amd64 go build -o gof-gen_osx gof-gen.go
CGO_ENABLED=0 GOOS=windows ARCH=amd64 go build -o gof-gen.exe gof-gen.go

tar cvzf generator_build.tar.gz gof-gen.sh gof-gen_* gof-gen.exe gen.conf code_templates
rm -rf gof-gen_* gof-gen.exe