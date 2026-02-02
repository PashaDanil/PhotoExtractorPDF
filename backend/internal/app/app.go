package app

import (
	"context"
	"imgpdf/internal/http/handler"
	"imgpdf/internal/minio"
	"imgpdf/internal/redis"
	"imgpdf/internal/service"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	server *http.Server
}

func New(ctx context.Context) (*App, error) {
	rdb, err := redis.New(ctx)
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(ctx)
	if err != nil {
		return nil, err
	}

	jobRedis := redis.NewJobRedis(rdb)
	minioStorage := minio.NewMinioStorage(minioClient)

	pdfService := service.NewPDFService()
	zipService := service.NewZIPService()
	jobService := service.NewJobService(jobRedis, minioStorage)

	PDFHandler := handler.NewPDFHandler(pdfService)
	ZIPHandler := handler.NewZIPHandler(zipService)
	jobHandler := handler.NewJobHandler(jobService)

	mux := http.NewServeMux()
	mux.HandleFunc("/pdf", PDFHandler.HandleTakePDF)
	mux.HandleFunc("/zip", ZIPHandler.HandleGiveZIP)
	mux.HandleFunc("/jobs", jobHandler.HandlePDFUploadRequest)
	mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jobs" && r.URL.Path != "/jobs/" {
			jobHandler.HandlePDFUploadComplete(w, r)
		} else {
			jobHandler.HandlePDFUploadRequest(w, r)
		}
	})
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("app initialized")

	return &App{
		server: server,
	}, nil
}

func (a *App) Run() error {

	log.Println("server started on :8080")

	if err := a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
