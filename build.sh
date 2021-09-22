#!/bin/bash

DEPLOYMENT_PACKAGE_NAME=deployment.zip
DEPLOYMENT_BUCKET_NAME=aws-serverless-poc
DEPLOYMENT_PREFIX_NAME=deployments
SAM_STACK_NAME=aws-serverless-poc
SAM_TEMPLATE_FILE=template.dev.yaml

echo "Building binaries..."
cd functions/get-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
rm main

cd ../create-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
rm main

cd ../search-books && GOOS=linux GOARCH=amd64 go build -o main cmd/main.go || exit
zip ${DEPLOYMENT_PACKAGE_NAME} main || exit
rm main

cd ../../sam || exit

echo "Building & packaging SAM template..."
sam build --template-file ./${SAM_TEMPLATE_FILE} || exit
sam package --s3-bucket ${DEPLOYMENT_BUCKET_NAME} --s3-prefix ${DEPLOYMENT_PREFIX_NAME} --template-file ./${SAM_TEMPLATE_FILE} --output-template-file ./template-output.dev.yaml || exit
echo "Deploying SAM template..."
sam deploy --template-file ./template-output.dev.yaml --stack-name ${SAM_STACK_NAME} --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND || exit

echo "Cleaning up..."
rm ../functions/get-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/create-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/search-books/${DEPLOYMENT_PACKAGE_NAME}

echo "Done."
