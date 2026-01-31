package service

type ZIPService struct{}

func NewZIPService() *ZIPService {
	return &ZIPService{}
}

func (s *ZIPService) GiveZIP() []byte {
	// берем архив с картинками из корня и отдаем его пользователю
	return []byte{}
}
