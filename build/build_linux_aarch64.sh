#!/bin/bash

sh build.sh

GOOS=linux GOARCH=arm64 go build -o ../aarch64 ../gouniversal.go