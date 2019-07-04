#!/bin/bash

FILE="version.go"
TIME=$(date +%Y%m%d_%H%M%S)
COMMIT=$(git rev-parse HEAD)

echo "package build\r\n" > $FILE
echo "const BuildTime = \""$TIME"\"" >> $FILE
echo "const Commit = \""$COMMIT"\"" >> $FILE
