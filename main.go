package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"io"
	"net/url"
	"strings"
)

func destinationUrl(input string) (string, error) {
	res, err := arn.Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse ARN %q: %w", input, err)
	}
	resourcePath, err := getResourceConsolePath(res)
	if err != nil {
		return "", err
	}

	if res.Region == "" {
		return fmt.Sprintf("https://console.aws.amazon.com%s", resourcePath), nil
	} else {
		return fmt.Sprintf("https://%s.console.aws.amazon.com%s", res.Region, resourcePath), nil
	}
}

func main() {
	reader, writer := io.Pipe()

	// Simulate writing to "stdin" by writing to the pipe
	go func() {
		defer writer.Close()

		arns := []string{
			"arn:aws:ecs:eu-west-1:0123456789:service/my-ecs-cluster/my-service",
			"arn:aws:sqs:eu-west-1:0123456789:job-dlq",
			"arn:aws:dynamodb:eu-west-1:0123456789:table/data-table",
			"arn:aws:s3:::some-bucket",
			"arn:aws:lambda:eu-west-1:0123456789:function:lambda-fn",
		}
		for _, a := range arns {
			a := a
			fmt.Fprintln(writer, a)
		}
	}()

	// Read from the pipe (which simulates stdin)
	scanner := bufio.NewScanner(reader)

	//scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Reading input (press Ctrl+D to stop):")
	for scanner.Scan() {
		line := scanner.Text()
		destinationUrl(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

func getResourceConsolePath(arn arn.ARN) (string, error) {
	fmt.Println(arn.Resource)
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
