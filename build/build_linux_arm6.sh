#!/bin/bash

GOOS=linux GOARCH=arm GOARM=6 go build -o ../arm6 ../gouniversal.go