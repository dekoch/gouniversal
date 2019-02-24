#!/bin/bash

sh build.sh

GOOS=linux GOARCH=arm GOARM=5 go build -o ../arm5 ../gouniversal.go