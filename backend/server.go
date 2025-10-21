package main

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	http *http.Server
}

func newServer(cfg *config) (*server, error) {
	presigner, err := initS3(cfg)
	if err != nil {
		return nil, err
	}

	svc := newService(presigner, cfg.AWS.Bucket, cfg.AWS.Folder, cfg.AWS.Region)
	hdl := newHandler(svc)

  r := httprouter.New()
	router(r, hdl)

	http := &http.Server{
		Addr:           ":5000",
		Handler:        r,
		MaxHeaderBytes: 100 * 1024 * 1024,
	}

	return &server{
		http,
	}, nil
}

func (s *server) start() error {
	return s.http.ListenAndServe()
}

func (s *server) shutdown(ctx context.Context) {
	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			log.Printf("Shutdown server thất bại: %v", err)
			return
		}
	}

	log.Println("Dừng server thành công")
}
