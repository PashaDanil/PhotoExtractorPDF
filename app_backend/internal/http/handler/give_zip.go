package handler

import (
	"imgpdf/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GiveZIPHandler struct {
	zipService *service.ZIPService
}

func NewGiveZIPHandler(zipService *service.ZIPService) *GiveZIPHandler {
	return &GiveZIPHandler{
		zipService: zipService,
	}
}

func (h *GiveZIPHandler) HandleGiveZIP(w http.ResponseWriter, r *http.Request) {
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
