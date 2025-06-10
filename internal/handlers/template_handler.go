package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
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

	safeFileName := filepath.Base(header.Filename)
	path := fmt.Sprintf("uploads/templates/company_%d_%s", companyID, safeFileName)

	// Создаём директорию, если нет
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "copy failed", http.StatusInternalServerError)
		return
	}

	t := &models.Template{
		CompanyID: companyID,
		Name:      name,
		FilePath:  path,
	}
	if err := h.Service.Create(t); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
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
