# PhotoExtractorPDF

Умная система для автоматического извлечения данных, таблиц и объектов из PDF документов с использованием компьютерного зрения и машинного обучения.

## Проблема

Система для обработки PDF-документов и извлечения изображений. Многие PDF-документы, в частности документации к технике, содержат векторные изображения и графические элементы в сложном внутреннем формате, из-за чего их невозможно корректно извлечь обычным парсингом структуры PDF. Проект решает задачу обнаружения и выделения данных графических элементов с использованием Computer Vision подходов.

PhotoExtractorPDF — это автоматизированный конвейер обработки PDF, который:
- Распознает структуру документов (текст, таблицы, изображения)
- Извлекает данные в структурированном виде
- Работает асинхронно с поддержкой очередей
- Масштабируется под любые объемы
- Предоставляет простой веб-интерфейс для загрузки и просмотра результатов


## Технологический стек

**Backend:** Go, REST API, RabbitMQ, Redis, MinIO

**Extractor:** Python, YOLOv11m-doclaynet, OpenCV, PyPDF

**Frontend:** React, TypeScript, Vite

**Инфраструктура:** Docker, Docker Compose

## Быстрый старт

### Требования
- Docker и Docker Compose
- или локально: Go 1.18+, Python 3.9+, Node.js 18+

### Запуск с Docker Compose

```bash
# Перейти в папку с конфигурацией
cd infra/compose

# Запустить все сервисы
docker-compose -f compose.yml up -d
```

### Доступ к сервисам
- **Frontend**: http://localhost:5173
- **API Swagger**: http://localhost:8080/swagger/index.html
- **MinIO Console**: http://localhost:9001

## Процесс обработки PDF

1. Пользователь загружает PDF через веб-интерфейс
2. Backend сохраняет файл в MinIO
3. Backend создает запись в Redis с состоянием задачи
4. Backend отправляет задачу в RabbitMQ
5. Extractor получает задачу из очереди
6. YOLOv11 распознает структуру документа
7. Извлекаются таблицы, текст, изображения
8. Результаты возвращаются в MinIO и Redis
9. Frontend запрашивает результаты из Redis
10. Пользователь видит извлеченные данные

## Разработка

### Backend (Go)
```bash
cd apps/backend
go mod download
go run cmd/api/main.go
```

### Extractor (Python)
```bash
cd apps/extractor
pip install -r requirements.txt
python main.py
```

### Frontend (React)
```bash
cd apps/frontend
npm install
npm run dev
```

## Лицензия

MIT
