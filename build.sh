#!/bin/bash

DEPLOYMENT_PACKAGE_NAME=deployment.zip

while getopts ":f:a" opt; do
  case $opt in
    f)
      echo "building binary for $OPTARG function..." >&2

      cd functions/"$OPTARG" && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      ;;
    a)
      echo "building binaries for all functions..." >&2

      cd functions/get-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../create-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../update-book && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../delete-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../search-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../stream-processor && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      cd ../auth && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
      zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
      rm main

      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

echo "Done."
