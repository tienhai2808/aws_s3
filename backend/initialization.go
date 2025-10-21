package main

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func initS3(cfg *config) (*s3.PresignClient, error) {
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

	client := s3.NewFromConfig(awsCfg)
	presigner := s3.NewPresignClient(client)
	
	return presigner, nil
}
