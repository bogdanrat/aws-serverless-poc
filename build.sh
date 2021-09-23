#!/bin/bash

DEPLOYMENT_REGION=eu-central-1
DEPLOYMENT_BUCKET_NAME=aws-serverless-poc
DEPLOYMENT_PREFIX_NAME=deployments
DEPLOYMENT_PACKAGE_NAME=deployment.zip
SAM_STACK_NAME=aws-serverless-poc
SAM_TEMPLATE_FILE=template.dev.yaml

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

cd ../../sam || exit

echo "Building & packaging SAM template..."
sam build --template-file ./${SAM_TEMPLATE_FILE} || exit
sam package --s3-bucket ${DEPLOYMENT_BUCKET_NAME} --s3-prefix ${DEPLOYMENT_PREFIX_NAME} --template-file ./${SAM_TEMPLATE_FILE} --output-template-file ./template-output.dev.yaml || exit
echo "Deploying SAM template..."
sam deploy --template-file ./template-output.dev.yaml --stack-name ${SAM_STACK_NAME} --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND --region ${DEPLOYMENT_REGION} --s3-bucket ${DEPLOYMENT_BUCKET_NAME} --parameter-overrides ContactEmailAddress=bogdanalexandru.rat@gmail.com || exit

apiName="$(aws cloudformation describe-stacks --stack-name ${SAM_STACK_NAME} --query "Stacks[0].Outputs[?OutputKey=='APIName'].OutputValue" --output text)"
# the AWS::ApiGateway::ApiKey resource created by the PER_API instruction has a logical ID of <api-logical-id>ApiKey
apiKeyId="$(aws cloudformation describe-stack-resources --stack-name ${SAM_STACK_NAME} --logical-resource-id "${apiName}"ApiKey --query "StackResources[0].PhysicalResourceId" --output text)"
# take the api key based on its physical resource id
apiKey="$(aws apigateway get-api-key --api-key "${apiKeyId}" --include-value)"

echo "API Key details: " "$apiKey"

echo "Cleaning up..."
rm ../functions/get-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/create-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/update-book/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/delete-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/search-books/${DEPLOYMENT_PACKAGE_NAME}
rm ../functions/stream-processor/${DEPLOYMENT_PACKAGE_NAME}

echo "Done."
