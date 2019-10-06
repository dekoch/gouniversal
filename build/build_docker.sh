#!/bin/bash

IMAGE="gou"
TAG="latest"

if [ "$1" != "" ]; then
	IMAGE=$1
fi

if [ "$2" != "" ]; then
	TAG=$2
fi

ERROR=0

for i in 1 2 3
do
    case $i in
        1)
            echo build application

            sh build.sh

            GOOS=linux go build -o ../gou ../gouniversal.go
            ERROR=$?
            ;;
        2)
            echo build docker image $IMAGE:$TAG

            docker build -t $IMAGE:$TAG ../.
            ERROR=$?
            ;;
        3)
            echo delete application file

            rm ../gou
            ERROR=$?
            ;;
    esac

    if [ "$ERROR" -ne "0" ]; then
        echo "error: $ERROR"
        break
    fi
done