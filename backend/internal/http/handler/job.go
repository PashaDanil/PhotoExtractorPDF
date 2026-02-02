package handler

import (
	"encoding/json"
	"imgpdf/internal/service"
	"log"
	"net/http"
	"strings"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

type InitUploadResponse struct {
	JobID        string `json:"job_id"`
	PDFKey       string `json:"pdf_key"`
	PresignedURL string `json:"presigned_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// HandlePDFUploadRequest godoc
// @Summary Инициализировать загрузку PDF
// @Description Создает новую задачу и возвращает presigned URL для загрузки PDF в MinIO
// @Tags jobs
// @Accept json
// @Produce json
// @Success 201 {object} InitUploadResponse "Данные для загрузки"
// @Failure 405 {object} ErrorResponse "Метод не разрешен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /jobs [post]
func (h *JobHandler) HandlePDFUploadRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	jobID, pdfKey, uploadURL, err := h.jobService.InitUpload(r.Context())
	if err != nil {
		log.Printf("init upload failed: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to initialize upload")
		return
	}

	// опционально: Location
	w.Header().Set("Location", "/jobs/"+jobID)

	writeJSON(w, http.StatusCreated, InitUploadResponse{
		JobID:        jobID,
		PDFKey:       pdfKey,
		PresignedURL: uploadURL,
	})
}

// HandlePDFUploadComplete godoc
// @Summary Завершить загрузку PDF
// @Description Отмечает задачу как готовую к обработке после успешной загрузки PDF
// @Tags jobs
// @Accept json
// @Produce json
// @Param job_id path string true "ID задачи"
// @Success 200 {object} map[string]string "Статус успешно обновлен"
// @Failure 400 {object} ErrorResponse "Некорректный путь или ID"
// @Failure 405 {object} ErrorResponse "Метод не разрешен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /jobs/{job_id}/complete-upload [post]
func (h *JobHandler) HandlePDFUploadComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	jobID, ok := parseCompleteUploadPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid path or job id")
		return
	}

	if err := h.jobService.CompleteUpload(r.Context(), jobID); err != nil {
		log.Printf("complete upload failed job_id=%s err=%v", jobID, err)
		// тут лучше маппить типы ошибок сервиса в 404/409 и т.п.
		writeError(w, http.StatusInternalServerError, "Failed to complete upload")
		return
	}

	// если дальше реально стартует async processing — я бы делал 202
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func parseCompleteUploadPath(p string) (string, bool) {
	// допускаем трейлинг-слэш
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	// ожидаем: ["jobs", "{job_id}", "complete-upload"]
	if len(parts) != 3 || parts[0] != "jobs" || parts[2] != "complete-upload" {
		return "", false
	}
	jobID := parts[1]
	if jobID == "" || strings.Contains(jobID, "/") {
		return "", false
	}
	return jobID, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
