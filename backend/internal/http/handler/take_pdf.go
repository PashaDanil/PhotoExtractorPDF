package handler

import (
	"imgpdf/internal/service"
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
// @Summary Загрузить PDF
// @Description Принимает PDF файл и сохраняет его для обработки
// @Tags pdf
// @Accept multipart/form-data
// @Produce json
// @Param pdf formData file true "PDF файл"
// @Success 200 "PDF успешно загружен"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /pdf [post]
func (h *PDFHandler) HandleTakePDF(w http.ResponseWriter, r *http.Request) {
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
