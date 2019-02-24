#!/bin/bash

sh build.sh

GOOS=linux GOARCH=arm GOARM=7 go build -o ../arm7 ../gouniversal.go