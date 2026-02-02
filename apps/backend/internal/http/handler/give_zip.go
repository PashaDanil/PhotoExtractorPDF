package handler

import (
	"go-api/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ZIPHandler struct {
	zipService *service.ZIPService
}

func NewZIPHandler(zipService *service.ZIPService) *ZIPHandler {
	return &ZIPHandler{
		zipService: zipService,
	}
}

// HandleGiveZIP godoc
// @ID getZip
// @Summary Download processed images as ZIP
// @Description Download all processed images from PDFs as a ZIP archive
// @Tags ZIP
// @Produce application/zip
// @Success 200 {file} file "ZIP file"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /zip [get]
func (h *ZIPHandler) HandleGiveZIP(w http.ResponseWriter, r *http.Request) {
	zip, err := h.zipService.GiveZIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		log.Printf("Error giving ZIP: %v", err)

		return
	}

	filename := "images-" + time.Now().Format("2006.01.02-15:04:05") + ".zip"

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(zip)))

	w.WriteHeader(http.StatusOK)
	w.Write(zip)

	log.Printf("ZIP file sent successfully: %s", filename)
}
