package handler

import (
	"go-api/internal/service"
	"io"
	"log"
	"net/http"
)

type PDFHandler struct {
	pdfService *service.PDFService
}

func NewPDFHandler(pdfService *service.PDFService) *PDFHandler {
	return &PDFHandler{
		pdfService: pdfService,
	}
}

// HandleTakePDF godoc
// @ID uploadPDF
// @Summary Upload PDF file
// @Description Upload a PDF file for processing
// @Tags PDF
// @Accept multipart/form-data
// @Produce json
// @Param pdf formData file true "PDF file to upload"
// @Success 200 {object} PDFUploadSuccessResponse "PDF uploaded successfully"
// @Failure 400 {object} ErrorResponse "Invalid request or file format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pdf [post]
func (h *PDFHandler) HandleTakePDF(w http.ResponseWriter, r *http.Request) {
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

	if err := h.pdfService.TakePDF(pdf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.Printf("Error saving PDF: %v", err)

		return
	}

	log.Println("PDF file saved successfully")

	writeJSON(w, http.StatusOK, PDFUploadSuccessResponse{
		Status:  "success",
		Message: "PDF uploaded successfully",
	})
}
