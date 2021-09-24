#!/bin/bash

DEPLOYMENT_REGION=eu-central-1
DEPLOYMENT_BUCKET_NAME=aws-serverless-poc
DEPLOYMENT_PREFIX_NAME=deployments
SAM_STACK_NAME=aws-serverless-poc
SAM_TEMPLATE_FILE=template.dev.yaml

cd sam || exit

echo "Building & packaging SAM template..."
sam build --template-file ./${SAM_TEMPLATE_FILE} || exit
sam package --s3-bucket ${DEPLOYMENT_BUCKET_NAME} --s3-prefix ${DEPLOYMENT_PREFIX_NAME} --region ${DEPLOYMENT_REGION} --template-file ./${SAM_TEMPLATE_FILE} --output-template-file ./template-output.dev.yaml || exit
echo "Deploying SAM template..."
sam deploy --template-file ./template-output.dev.yaml --stack-name ${SAM_STACK_NAME} --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND --region ${DEPLOYMENT_REGION} --s3-bucket ${DEPLOYMENT_BUCKET_NAME} --parameter-overrides ContactEmailAddress=bogdanalexandru.rat@gmail.com || exit

apiName="$(aws cloudformation describe-stacks --stack-name ${SAM_STACK_NAME} --region ${DEPLOYMENT_REGION} --query "Stacks[0].Outputs[?OutputKey=='APIName'].OutputValue" --output text)"
# the AWS::ApiGateway::ApiKey resource created by the PER_API instruction has a logical ID of <api-logical-id>ApiKey
apiKeyId="$(aws cloudformation describe-stack-resources --stack-name ${SAM_STACK_NAME} --region ${DEPLOYMENT_REGION} --logical-resource-id "${apiName}"ApiKey --query "StackResources[0].PhysicalResourceId" --output text)"
# take the api key based on its physical resource id
apiKey="$(aws apigateway get-api-key --region ${DEPLOYMENT_REGION} --api-key "${apiKeyId}" --include-value)"

echo "API Key details: " "$apiKey"

echo "Done."
