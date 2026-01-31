package app

import (
	"imgpdf/internal/http/handler"
	"imgpdf/internal/service"
	"log"
	"net/http"
)

type App struct {
	server *http.Server
}

func New() (*App, error) {
	pdfService := service.NewPDFService()
	zipService := service.NewZIPService()

	TakePDFHandler := handler.NewTakePDFHandler(pdfService)
	GiveZIPHandler := handler.NewGiveZIPHandler(zipService)

	mux := http.NewServeMux()
	mux.HandleFunc("/pdf", TakePDFHandler.HandleTakePDF)
	mux.HandleFunc("/zip", GiveZIPHandler.HandleGiveZIP)

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
