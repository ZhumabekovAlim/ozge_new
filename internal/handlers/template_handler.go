package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type TemplateHandler struct {
	Service *services.TemplateService
}

func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{Service: service}
}

func (h *TemplateHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "could not get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	companyID, _ := strconv.Atoi(r.FormValue("company_id"))
	name := r.FormValue("name")

	// --- Создаём директорию для компании
	companyDir := fmt.Sprintf("C:\\Users\\alimz\\GolandProjects\\OzgeContract\\uploads\\templates\\company_%d", companyID)
	if err := os.MkdirAll(companyDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}
	safeFileName := filepath.Base(header.Filename)
	origPath := filepath.Join(companyDir, safeFileName)
	watermarkedFileName := fmt.Sprintf("watermarked_%s", safeFileName)
	watermarkedPath := filepath.Join(companyDir, watermarkedFileName)

	// --- Сохраняем оригинал
	tmpFile, err := os.Create(origPath)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(origPath)
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "copy failed", http.StatusInternalServerError)
		return
	}

	// --- Watermark
	desc := "op:0.2, pos:bottom-right, scale:0.1, rot:0, off:-20 20"
	wm, err := pdfcpu.ParseImageWatermarkDetails("C:\\Users\\alimz\\GolandProjects\\OzgeContract\\static\\contract-logo.png", desc, false, types.POINTS)
	if err != nil {
		http.Error(w, "cannot parse image watermark: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.AddWatermarksFile(origPath, watermarkedPath, []string{"1-"}, wm, nil)
	if err != nil {
		http.Error(w, "watermark error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t := &models.Template{
		CompanyID: companyID,
		Name:      name,
		FilePath:  watermarkedPath,
	}

	if err := h.Service.Create(t); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	t, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(t)
	w.WriteHeader(http.StatusOK)
}

func (h *TemplateHandler) GetByCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	list, err := h.Service.GetByCompany(id)
	if err != nil {
		http.Error(w, "cannot fetch templates", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
	w.WriteHeader(http.StatusOK)
}

func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)

	var input models.Template
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	input.ID = id

	if err := h.Service.Update(&input); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TemplateHandler) ServePDFByID(w http.ResponseWriter, r *http.Request) {
	// Получаем id шаблона
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Получаем шаблон из БД
	template, err := h.Service.GetByID(id)
	if err != nil || template == nil {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}

	// Безопасность: только basename файла
	filePath := template.FilePath
	// Если хранится с обратными слэшами (\), заменяем на /
	filePath = filepath.FromSlash(filePath)

	// Проверяем, что файл существует
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// Заголовки для отображения PDF в браузере
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(filePath)+"\"")
	http.ServeFile(w, r, filePath)
}
