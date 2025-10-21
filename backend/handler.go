package main

import (
	"context"
	"encoding/json"
	"net/http"
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
	toResponse(w, http.StatusOK, "Welcome to the home page", nil)
}

func (h *handler) generatePresignedURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req presignedURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		toResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if req.FileName == "" || req.ContentType == "" {
		toResponse(w, http.StatusBadRequest, "fileName và contentType là bắt buộc", nil)
		return
	}

	result, err := h.svc.CreateUploadURL(ctx, req)
	if err != nil {
		toResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	toResponse(w, http.StatusOK, "Tạo presigned URL thành công", result)
}
