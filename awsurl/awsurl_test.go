package awsurl

import (
	"testing"
)

func TestDestinationUrl(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ECS ARN",
			input:    "arn:aws:ecs:eu-west-1:0123456789:service/my-ecs-cluster/my-service",
			expected: "https://eu-west-1.console.aws.amazon.com/ecs/v2/clusters/my-ecs-cluster/services/my-service",
		},
		{
			name:     "SQS ARN",
			input:    "arn:aws:sqs:eu-west-1:0123456789:job-dlq",
			expected: "https://eu-west-1.console.aws.amazon.com/sqs/v3/home#/queues/https%3A%2F%2Fsqs.eu-west-1.amazonaws.com%2F0123456789%2Fjob-dlq",
		},
		{
			name:     "DynamoDB ARN",
			input:    "arn:aws:dynamodb:eu-west-1:0123456789:table/data-table",
			expected: "https://eu-west-1.console.aws.amazon.com/dynamodbv2/home#table?name=table/data-table",
		},
		{
			name:     "S3 ARN",
			input:    "arn:aws:s3:::some-bucket",
			expected: "https://console.aws.amazon.com/s3/buckets/some-bucket",
		},
		{
			name:     "Lambda ARN",
			input:    "arn:aws:lambda:eu-west-1:0123456789:function:lambda-fn",
			expected: "https://eu-west-1.console.aws.amazon.com/lambda/home#/functions/lambda-fn?accountID=0123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromArn(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.expected {
				t.Errorf("FromArn() = %v, want %v", got, tt.expected)
			}
		})
	}
}
