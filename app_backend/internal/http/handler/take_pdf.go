package handler

import "imgpdf/internal/service"

type TakePDFHandler struct {
	pdfService *service.PDFService
}

func NewTakePDFHandler(pdfService *service.PDFService) *TakePDFHandler {
	return &TakePDFHandler{
		pdfService: pdfService,
	}
}

func (h *TakePDFHandler) HandleTakePDF() {
	// берем pdf из запроса
	h.pdfService.TakePDF( /* pdf документ */ )
}
