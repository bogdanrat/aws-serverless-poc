package cwlogger

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/bogdanrat/aws-serverless-poc/contracts/common"
	"github.com/bogdanrat/aws-serverless-poc/lib/logger"
)

type CWMetricLogger struct {
	Client *cloudwatch.Client
}

func New(cfg aws.Config) logger.MetricLogger {
	return &CWMetricLogger{
		Client: cloudwatch.NewFromConfig(cfg),
	}
}

func (l *CWMetricLogger) PutMetric(dimensionName string, metricName string) error {
	_, err := l.Client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(common.BooksMetricNamespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Value:      aws.Float64(1.0),
				Unit:       types.StandardUnitCount,
				Dimensions: []types.Dimension{
					{
						Name:  aws.String(common.EnvironmentDimensionName),
						Value: aws.String(dimensionName),
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}
	return nil
}
