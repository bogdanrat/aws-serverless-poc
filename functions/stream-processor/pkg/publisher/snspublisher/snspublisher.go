package snspublisher

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/functions/stream-processor/pkg/publisher"
	"os"
)

type SNSPublisher struct {
	Client   *sns.Client
	TopicARN string
}

func New(cfg aws.Config) publisher.Publisher {
	return &SNSPublisher{
		Client:   sns.NewFromConfig(cfg),
		TopicARN: os.Getenv(common.BooksTopicArnEnvironmentVariable),
	}
}

func (p *SNSPublisher) Publish(message string) error {
	input := &sns.PublishInput{
		Message:   aws.String(message),
		TargetArn: aws.String(p.TopicARN),
	}

	_, err := p.Client.Publish(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
