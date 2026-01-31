package service

import (
	"fmt"
	"os"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) TakePDF(pdf []byte) error {
	// берем pdf и загружаем его в корень
	if len(pdf) == 0 {
		return fmt.Errorf("pdf is empty")
	}
	if err := os.WriteFile("file.pdf", pdf, 0644); err != nil {
		return fmt.Errorf("pdf downloading error: %s", err)
	}
	return nil
}
