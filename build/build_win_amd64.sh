#!/bin/bash

sh build.sh

GOOS=windows GOARCH=amd64 go build -o ../gouniversal_amd64.exe ../gouniversal.go