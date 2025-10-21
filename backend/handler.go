package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	svc service
}

func newHandler(svc service) *handler {
	return &handler{svc}
}

func (h *handler) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	toApiResponse(w, http.StatusOK, "Welcome to the home page", nil)
}

func (h *handler) generatePresignedURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req presignedURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		toApiResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if req.FileName == "" || req.ContentType == "" {
		toApiResponse(w, http.StatusBadRequest, "fileName và contentType là bắt buộc", nil)
		return
	}

	result, err := h.svc.createUploadURL(ctx, req)
	if err != nil {
		toApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	toApiResponse(w, http.StatusOK, "Tạo presigned URL upload file thành công", result)
}

func (h *handler) viewFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	objectKey := r.URL.Query().Get("key")
	if objectKey == "" {
		toApiResponse(w, http.StatusBadRequest, "tham số truy vấn 'key' là bắt buộc", nil)
		return
	}

	fileUrl, err := h.svc.createViewURL(ctx, objectKey)
	if err != nil {
		switch err {
		case errFileNotFound:
			toApiResponse(w, http.StatusNotFound, err.Error(), nil)
		default:
			toApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	toApiResponse(w, http.StatusOK, "Tạo presigned URL xem file thành công", map[string]string{
		"url": fileUrl,
	})
}

func (h *handler) deleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rawKey := ps.ByName("key")
	objectKey := strings.TrimPrefix(rawKey, "/")
	if objectKey == "" {
		toApiResponse(w, http.StatusBadRequest, "key trong URL là bắt buộc", nil)
		return
	}

	if err := h.svc.deleteFile(ctx, objectKey); err != nil {
		switch err {
		case errFileNotFound:
			toApiResponse(w, http.StatusNotFound, err.Error(), nil)
		default:
			toApiResponse(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	toApiResponse(w, http.StatusOK, "Xóa file trên S3 thành công", nil)
}
