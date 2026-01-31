package handler

import (
	"imgpdf/internal/service"
	"io"
	"log"
	"net/http"
)

type TakePDFHandler struct {
	pdfService *service.PDFService
}

func NewTakePDFHandler(pdfService *service.PDFService) *TakePDFHandler {
	return &TakePDFHandler{
		pdfService: pdfService,
	}
}

func (h *TakePDFHandler) HandleTakePDF(w http.ResponseWriter, r *http.Request) {
	// берем pdf из запроса
	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.Printf("Error reading PDF from request: %v", err)

		return
	}
	defer file.Close()

	pdf, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		log.Printf("Error reading PDF data: %v", err)

		return
	}

	// вызываем сервис для сохранения pdf
	if err := h.pdfService.TakePDF(pdf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.Printf("Error saving PDF: %v", err)

		return
	}

	log.Println("PDF file saved successfully")

	w.WriteHeader(http.StatusOK)
}
