package handler

import "imgpdf/internal/service"

type GiveZIPHandler struct {
	zipService *service.ZIPService
}

func NewGiveZIPHandler(zipService *service.ZIPService) *GiveZIPHandler {
	return &GiveZIPHandler{
		zipService: zipService,
	}
}

func (h *GiveZIPHandler) HandleGiveZIP() {
	zip := h.zipService.GiveZIP( /* имя файла */ )
	_ = zip
}
