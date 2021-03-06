# AWS Serverless CRUD API POC
## Managed, built & deployed through AWS SAM

### Services
- **Amazon API Gateway**: proxies incoming request to the corresponding Lambda function, as an AWS_PROXY integration. Integrates with **Amazon Cognito** for Authorization. Features an API Usage Plan for API Quota & Throttling and API Key for Client access.
- **AWS Lambda**: each Lambda function handles a specific event: create, get, update, delete, search (by book author, title or category) and auth events (signup and login)
- **Amazon DynamoDB**: models a Book entity, with Author & Title as the primary key (hash key and sort key, respectively). Item-level changes are pushed to a **DynamoDB Streams**. An **Event Source Mapping** is established between the streams and a Lambda function that processes the stream and sends notifications to an **SNS Topic** when:
  - a new book has been published
  - a book went out of stock
  - new formats (hardcover, paperback, audiobook etc.) are available for a book
- **Amazon SNS**: implements a Pub/Sub pattern for book status changes. The email subscription is automatically created through the SAM template by the stack creation responsible (configurable through parameters).
- **Amazon Cognito**: provides authentication and authorization mechanisms
- **Amazon CloudWatch**: monitoring custom metric data sent by the Lambda functions

![Architecture](https://bogdans-public-assets.s3.eu-central-1.amazonaws.com/AWS%20API%20POC.png)
