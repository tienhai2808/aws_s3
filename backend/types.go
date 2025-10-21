package main

import "context"

type service interface {
	createUploadURL(ctx context.Context, req presignedURLRequest) (*presignedURLResponse, error)

	createViewURL(ctx context.Context, key string) (string, error)

	deleteFile(ctx context.Context, key string) error
}

type presignedURLRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type presignedURLResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

type apiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
