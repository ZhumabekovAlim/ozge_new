package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type ContractHandler struct {
	Service *services.ContractService
}

func NewContractHandler(service *services.ContractService) *ContractHandler {
	return &ContractHandler{Service: service}
}

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	companyID, _ := strconv.Atoi(r.FormValue("company_id"))
	templateID, _ := strconv.Atoi(r.FormValue("template_id"))
	method := r.FormValue("method")
	clientFilled := r.FormValue("client_filled") == "true"
	companySign := r.FormValue("company_sign") == "true" || r.FormValue("company_sign") == "1"

	input := models.Contract{
		CompanyID:     companyID,
		TemplateID:    templateID,
		ClientFilled:  clientFilled,
		Method:        method,
		CompanySign:   companySign,
		ContractToken: uuid.New().String(),
	}

	if err := h.Service.Create(&input); err != nil {
		http.Error(w, "failed to create contract", http.StatusInternalServerError)
		return
	}

	baseDir := os.Getenv("DATA_DIR")
	if baseDir == "" {
		baseDir = "uploads"
	}
	contractsDir := filepath.Join(baseDir, "contracts", fmt.Sprintf("company_%d", input.CompanyID))
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	finalPDF := filepath.Join(contractsDir, fmt.Sprintf("final_%d.pdf", input.ID))
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	out, err := os.Create(finalPDF)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}

	input.GeneratedPDFPath = finalPDF
	_ = h.Service.UpdatePDFPath(input.ID, finalPDF)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

func (h *ContractHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	contract, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(contract)
}

func (h *ContractHandler) GetByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")
	contract, err := h.Service.GetByToken(token)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(contract)
}

func (h *ContractHandler) GetByTokenWithFields(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")
	details, err := h.Service.GetByTokenWithFields(token)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(details)
}

func (h *ContractHandler) GetByCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	list, err := h.Service.GetByCompanyID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	var input models.Contract
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

func (h *ContractHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ContractHandler) CreateWithFields(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	companyID, _ := strconv.Atoi(r.FormValue("company_id"))
	templateID, _ := strconv.Atoi(r.FormValue("template_id"))
	method := r.FormValue("method")
	clientFilled := r.FormValue("client_filled") == "true"
	companySign := r.FormValue("company_sign") == "true" || r.FormValue("company_sign") == "1"

	var fieldDTOs []models.ContractFieldDTO
	if val := r.FormValue("fields"); val != "" {
		_ = json.Unmarshal([]byte(val), &fieldDTOs)
	}

	contract := models.Contract{
		CompanyID:     companyID,
		TemplateID:    templateID,
		ClientFilled:  clientFilled,
		Method:        method,
		CompanySign:   companySign,
		ContractToken: uuid.New().String(),
	}

	fields := make([]models.ContractField, 0, len(fieldDTOs))
	for _, f := range fieldDTOs {
		fields = append(fields, models.ContractField{
			FieldName: f.FieldName,
			FieldType: f.FieldType,
		})
	}

	if err := h.Service.CreateWithFields(&contract, fields); err != nil {
		http.Error(w, "failed to create contract and fields", http.StatusInternalServerError)
		return
	}

	baseDir := os.Getenv("DATA_DIR")
	if baseDir == "" {
		baseDir = "uploads"
	}
	contractsDir := filepath.Join(baseDir, "contracts", fmt.Sprintf("company_%d", contract.CompanyID))
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	finalPDF := filepath.Join(contractsDir, fmt.Sprintf("final_%d.pdf", contract.ID))
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	out, err := os.Create(finalPDF)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}

	contract.GeneratedPDFPath = finalPDF
	_ = h.Service.UpdatePDFPath(contract.ID, finalPDF)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contract)
}

// ServePDFByID streams the generated contract PDF to the client.
func (h *ContractHandler) ServePDFByID(w http.ResponseWriter, r *http.Request) {
	// Extract contract id from the URL
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Retrieve contract information from the database
	contract, err := h.Service.GetByID(id)
	if err != nil || contract == nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}

	// Normalize path and ensure the file exists
	filePath := filepath.FromSlash(contract.GeneratedPDFPath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// Serve the PDF inline in the browser
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(filePath)+"\"")
	http.ServeFile(w, r, filePath)
}
