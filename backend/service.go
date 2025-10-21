package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type serviceImpl struct {
	presigner *s3.PresignClient
	bucket    string
	folder    string
	region    string
}

func newService(presigner *s3.PresignClient, bucket, folder, region string) service {
	return &serviceImpl{
		presigner,
		bucket,
		folder,
		region,
	}
}

func (s *serviceImpl) CreateUploadURL(ctx context.Context, req presignedURLRequest) (*presignedURLResponse, error) {
	objectKey := fmt.Sprintf("%s/%s-%s", s.folder, uuid.New().String(), req.FileName)

	presignedReq, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(req.ContentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return nil, fmt.Errorf("generate presigned URL thất bại: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, objectKey)

	return &presignedURLResponse{
		UploadURL: presignedReq.URL,
		FileURL: fileURL,
	}, nil
}
