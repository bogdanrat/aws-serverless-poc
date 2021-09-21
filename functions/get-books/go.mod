module github.com/bogdanrat/aws-serverless-poc/functions/get-books

go 1.16

require (
	github.com/aws/aws-lambda-go v1.26.0
	github.com/aws/aws-sdk-go-v2 v1.9.1
	github.com/aws/aws-sdk-go-v2/config v1.8.2
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.2.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.8.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.5.1
)
