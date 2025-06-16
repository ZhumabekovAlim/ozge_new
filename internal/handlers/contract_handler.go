package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ContractHandler struct {
	Service *services.ContractService
}

func NewContractHandler(service *services.ContractService) *ContractHandler {
	return &ContractHandler{Service: service}
}

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.Contract
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	input.ContractToken = uuid.New().String()

	// 1. Сохраняем контракт в базе (input.ID заполняется после Create)
	err := h.Service.Create(&input)
	if err != nil {
		http.Error(w, "failed to create contract", http.StatusInternalServerError)
		return
	}

	// 2. Создаём директорию для компании
	contractsDir := fmt.Sprintf("uploads/contracts/company_%d", input.CompanyID)
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	// 3. Пути к файлам
	mainPDF := input.GeneratedPDFPath // исходный PDF с ватермаркой
	certPDF := filepath.Join(contractsDir, fmt.Sprintf("cert_%d.pdf", input.ID))
	finalPDF := filepath.Join(contractsDir, fmt.Sprintf("final_%d.pdf", input.ID))

	// 4. Имя компании (можешь получить из базы)
	companyName := "Test Company"

	// 5. Генерим страницу сертификата
	if err := GenerateCertificatePDF(certPDF, companyName, time.Now()); err != nil {
		http.Error(w, "certificate page error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Мержим основной PDF и сертификат
	if err := api.MergeCreateFile([]string{mainPDF, certPDF}, finalPDF, false, nil); err != nil {
		http.Error(w, "merge error: "+err.Error(), http.StatusInternalServerError)
		_ = os.Remove(certPDF)
		return
	}

	// 7. Удаляем временный сертификат
	_ = os.Remove(certPDF)

	// 8. Обновляем путь финального файла в контракте
	input.GeneratedPDFPath = finalPDF
	_ = h.Service.UpdatePDFPath(input.ID, finalPDF)

	// 9. Вернём контракт с финальным путём
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
	var input models.CreateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	contract := models.Contract{
		CompanyID:        input.CompanyID,
		TemplateID:       input.TemplateID,
		GeneratedPDFPath: input.GeneratedPDFPath,
		ClientFilled:     input.ClientFilled,
		Method:           input.Method,
		ContractToken:    uuid.New().String(),
	}

	fields := make([]models.ContractField, 0)
	for _, f := range input.Fields {
		fields = append(fields, models.ContractField{
			FieldName: f.FieldName,
			FieldType: f.FieldType,
		})
	}

	// Создаём контракт и поля, contract.ID будет получен после сохранения
	err := h.Service.CreateWithFields(&contract, fields)
	if err != nil {
		http.Error(w, "failed to create contract and fields", http.StatusInternalServerError)
		return
	}

	// --- Создаём директорию для компании
	contractsDir := fmt.Sprintf("uploads/contracts/company_%d", contract.CompanyID)
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	mainPDF := contract.GeneratedPDFPath
	certPDF := filepath.Join(contractsDir, fmt.Sprintf("cert_%d.pdf", contract.ID))
	finalPDF := filepath.Join(contractsDir, fmt.Sprintf("final_%d.pdf", contract.ID))

	companyName := "Test Company"

	if err := GenerateCertificatePDF(certPDF, companyName, time.Now()); err != nil {
		http.Error(w, "certificate page error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := api.MergeCreateFile([]string{mainPDF, certPDF}, finalPDF, false, nil); err != nil {
		http.Error(w, "merge error: "+err.Error(), http.StatusInternalServerError)
		_ = os.Remove(certPDF)
		return
	}

	_ = os.Remove(certPDF)

	contract.GeneratedPDFPath = finalPDF
	_ = h.Service.UpdatePDFPath(contract.ID, finalPDF)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contract)
}

func GenerateCertificatePDF(filename, companyName string, signedAt time.Time) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 20)
	pdf.AddPage()
	pdf.CellFormat(0, 20, "Сертификат онлайн подписания", "", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("DejaVu", "", 14)
	pdf.SetXY(10, 40)
	pdf.MultiCell(0, 10, fmt.Sprintf(
		"Сторона 1\nИсполнитель: %s\nПодписан: %s",
		companyName,
		signedAt.Format("2006-01-02 15:04:05"),
	), "", "L", false)
	return pdf.OutputFileAndClose(filename)
}
