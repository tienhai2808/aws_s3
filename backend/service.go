package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type serviceImpl struct {
	client    *awsS3.Client
	presigner *awsS3.PresignClient
	bucket    string
	folder    string
	region    string
}

func newService(client *awsS3.Client, presigner *awsS3.PresignClient, bucket, folder, region string) service {
	return &serviceImpl{
		client,
		presigner,
		bucket,
		folder,
		region,
	}
}

func (s *serviceImpl) createUploadURL(ctx context.Context, req presignedURLRequest) (*presignedURLResponse, error) {
	objectKey := fmt.Sprintf("%s/%s-%s", s.folder, uuid.New().String(), req.FileName)

	presignedReq, err := s.presigner.PresignPutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(req.ContentType),
	}, func(opts *awsS3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return nil, fmt.Errorf("generate presigned URL upload file thất bại: %w", err)
	}

	return &presignedURLResponse{
		UploadURL: presignedReq.URL,
		ObjectKey: objectKey,
	}, nil
}

func (s *serviceImpl) createViewURL(ctx context.Context, key string) (string, error) {
	if _, err := s.client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		var keyNotFound *types.NotFound
		if errors.As(err, &keyNotFound) {
			return "", errFileNotFound
		}
		return "", fmt.Errorf("kiểm tra file thất bại: %w", err)
	}

	presignedReq, err := s.presigner.PresignGetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *awsS3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})
	if err != nil {
		return "", fmt.Errorf("generate presigned URL xem file thất bại: %w", err)
	}

	return presignedReq.URL, nil
}

func (s *serviceImpl) deleteFile(ctx context.Context, key string) error {
	if _, err := s.client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		var keyNotFound *types.NotFound
		if errors.As(err, &keyNotFound) {
			return errFileNotFound
		}
		return fmt.Errorf("kiểm tra file thất bại: %w", err)
	}

	if _, err := s.client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("xóa file trên S3 thất bại: %w", err)
	}

	return nil
}
