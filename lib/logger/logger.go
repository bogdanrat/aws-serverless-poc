package logger

type MetricLogger interface {
	PutMetric(dimensionName string, metricName string) error
}
