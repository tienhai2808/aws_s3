package main

import "context"

type service interface {
	CreateUploadURL(ctx context.Context, req presignedURLRequest) (*presignedURLResponse, error)
}

type presignedURLRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type presignedURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
}

type apiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
