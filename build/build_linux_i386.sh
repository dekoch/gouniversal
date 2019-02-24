#!/bin/bash

sh build.sh

GOOS=linux GOARCH=386 go build -o ../i386 ../gouniversal.go