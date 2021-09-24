#!/bin/bash

DEPLOYMENT_PACKAGE_NAME=deployment.zip

echo "Building binaries..."
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

echo "Done."
