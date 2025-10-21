package main

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3 struct {
	client    *awsS3.Client
	presigner *awsS3.PresignClient
}

func initS3(cfg *config) (*s3, error) {
	staticCredentials := credentials.NewStaticCredentialsProvider(
		cfg.AWS.AccessKeyID,
		cfg.AWS.SecretAccessKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(cfg.AWS.Region),
		awsConfig.WithCredentialsProvider(staticCredentials),
	)
	if err != nil {
		return nil, fmt.Errorf("tải cấu hình aws thất bại")
	}

	client := awsS3.NewFromConfig(awsCfg)
	presigner := awsS3.NewPresignClient(client)

	return &s3{
		client,
		presigner,
	}, nil
}
