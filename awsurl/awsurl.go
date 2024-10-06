package awsurl

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"net/url"
	"strings"
)

func pathFromArn(arn arn.ARN) (string, error) {
	switch arn.Service {
	case "ecs":
		parts := strings.Split(arn.Resource, "/")
		if len(parts) < 3 || parts[0] != "service" {
			return "", fmt.Errorf("unexpected ARN format for ECS: %s", arn.Resource)
		}
		cluster := parts[1]
		service := parts[2]
		return fmt.Sprintf("/ecs/v2/clusters/%s/services/%s", cluster, service), nil

	case "dynamodb":
		tableName := arn.Resource
		return fmt.Sprintf("/dynamodbv2/home#table?name=%s", tableName), nil

	case "s3":
		bucketName := arn.Resource
		return fmt.Sprintf("/s3/buckets/%s", bucketName), nil

	case "lambda":
		parts := strings.Split(arn.Resource, ":")
		if len(parts) < 2 || parts[0] != "function" {
			return "", fmt.Errorf("unexpected ARN format for ECS: %s", arn.Resource)
		}

		functionName := parts[1]
		return fmt.Sprintf("/lambda/home#/functions/%s?accountID=%s", functionName, arn.AccountID), nil

	case "sqs":
		sqsUrl := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", arn.Region, arn.AccountID, arn.Resource)
		encodedUrl := url.QueryEscape(sqsUrl)
		return fmt.Sprintf("/sqs/v3/home#/queues/%s", encodedUrl), nil
	default:
		return "", fmt.Errorf("unknown service in ARN: %s", arn.Resource)
	}
}

func urlWithOptionalRegion(res arn.ARN, path string) string {
	if res.Region == "" {
		return fmt.Sprintf("https://console.aws.amazon.com%s", path)
	} else {
		return fmt.Sprintf("https://%s.console.aws.amazon.com%s", res.Region, path)
	}
}

func FromArn(input string) (string, error) {
	res, err := arn.Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse ARN %q: %w", input, err)
	}
	path, err := pathFromArn(res)
	if err != nil {
		return "", err
	}

	return urlWithOptionalRegion(res, path), nil
}

type Resolver struct {
	AccessPortalDomain string
	RoleName           string
}

func (c *Resolver) FromArn2(input string) (string, error) {
	res, err := arn.Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse ARN %q: %w", input, err)
	}
	path, err := pathFromArn(res)
	if err != nil {
		return "", err
	}
	awsUrl := urlWithOptionalRegion(res, path)

	if c.AccessPortalDomain != "" {
		return fmt.Sprintf("https://%s.awsapps.com/start/#/console?account_id=%s&role_name=%s&destination=%s",
			c.AccessPortalDomain,
			res.AccountID,
			c.RoleName,
			awsUrl,
		), nil
	} else {
		return awsUrl, nil
	}

}
