#!/bin/bash

GOOS=linux GOARCH=arm GOARM=7 go build -o ../arm7 ../gouniversal.go
