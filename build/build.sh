#!/bin/bash

FILE="version.go"
TIME=$(date +%Y%m%d_%H%M%S)

echo "package build\r\n" > $FILE
echo "const BuildTime = \""$TIME"\"" >> $FILE