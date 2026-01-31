package service

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) TakePDF() {
	// берем pdf из запроса и загружаем его в корень

}
