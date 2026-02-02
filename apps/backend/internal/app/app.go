package app

import (
	"context"
	"go-api/internal/http/handler"
	"go-api/internal/redis"
	"go-api/internal/service"
	"go-api/internal/storage"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	// httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	server *http.Server
}

func New(ctx context.Context) (*App, error) {
	_ = godotenv.Load(".env")

	rdb, err := redis.New(ctx)
	if err != nil {
		return nil, err
	}

	minioClient, err := storage.New(ctx)
	if err != nil {
		return nil, err
	}

	jobRedis := redis.NewJobRedis(rdb)
	minioStorage := storage.NewStorage(minioClient)

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
	// mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

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
