#!/bin/bash

sh build.sh

GOOS=linux GOARCH=amd64 go build -o ../amd64 ../gouniversal.go