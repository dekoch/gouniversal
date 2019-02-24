#!/bin/bash

sh build.sh

GOOS=linux GOARCH=arm GOARM=6 go build -o ../arm6 ../gouniversal.go