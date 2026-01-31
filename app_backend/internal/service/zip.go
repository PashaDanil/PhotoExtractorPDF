package service

import (
	"fmt"
	"os"
)

type ZIPService struct{}

func NewZIPService() *ZIPService {
	return &ZIPService{}
}

func (s *ZIPService) GiveZIP() ([]byte, error) {
	// берем архив с картинками из корня и отдаем его пользователю
	zip, err := os.ReadFile("images.zip")
	if err != nil {
		return nil, fmt.Errorf("zip reading error: %w", err)
	}
	return zip, nil
}
