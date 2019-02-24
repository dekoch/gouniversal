#!/bin/bash

sh build.sh

GOOS=windows GOARCH=386 go build -o ../gouniversal_i386.exe ../gouniversal.go