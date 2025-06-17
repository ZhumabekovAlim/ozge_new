package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	_ "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type SignatureHandler struct {
	Service *services.SignatureService
}

func NewSignatureHandler(service *services.SignatureService) *SignatureHandler {
	return &SignatureHandler{Service: service}
}

func (h *SignatureHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.Signature
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Создание в базе
	newID, err := h.Service.Create(&input)
	if err != nil {
		http.Error(w, "failed to create signature", http.StatusInternalServerError)
		return
	}
	input.ID = newID

	// Получение контракта
	contract, err := h.Service.GetContractByID(input.ContractID)
	if err != nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}
	contractPDF := contract.GeneratedPDFPath

	// Каталог для файлов
	signDir := fmt.Sprintf("uploads/signatures/company_%d", contract.CompanyID)
	if err := os.MkdirAll(signDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	// Генерация QR
	signedPDF := filepath.Join(signDir, fmt.Sprintf("signed_final_%d.pdf", input.ID))
	qrPath := filepath.Join(signDir, fmt.Sprintf("qr_%d.png", input.ID))

	// ⬇️ ссылка на просмотр по ID
	qrContent := fmt.Sprintf("http://192.168.8.2:4000/signatures/pdf/%d", input.ID)

	err = qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrPath)
	if err != nil {
		http.Error(w, "QR generation error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(qrPath)
	err = qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrPath)
	if err != nil {
		http.Error(w, "QR generation error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(qrPath)

	// Добавление QR в PDF
	err = AddQRToPDF(contractPDF, signedPDF, qrPath)
	if err != nil {
		http.Error(w, "failed to insert QR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновление пути
	input.SignFilePath = signedPDF
	_ = h.Service.UpdateSignFilePath(input.ID, signedPDF)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func AddQRToPDF(inputPDF, outputPDF, qrImagePath string) error {
	conf := model.NewDefaultConfiguration()

	// 1. Узнаем количество страниц в PDF
	ctx, err := api.ReadContextFile(inputPDF)
	if err != nil {
		return fmt.Errorf("failed to read PDF context: %w", err)
	}
	totalPages := ctx.PageCount
	lastPage := fmt.Sprintf("%d", totalPages)

	// 2. Подготавливаем QR watermark
	wm, err := pdfcpu.ParseImageWatermarkDetails(
		qrImagePath,
		"pos:tr, scale:0.2, rot:0, offset:-10 -100",
		true, // onTop
		types.POINTS,
	)
	if err != nil {
		return fmt.Errorf("watermark parse error: %w", err)
	}

	// 3. Добавляем на последнюю страницу
	return api.AddWatermarksFile(inputPDF, outputPDF, []string{lastPage}, wm, conf)
}

func (h *SignatureHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	sig, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusOK)
		return
	}
	json.NewEncoder(w).Encode(sig)
}

func (h *SignatureHandler) GetByContractID(w http.ResponseWriter, r *http.Request) {
	contractIDStr := r.URL.Query().Get(":id")
	contractID, _ := strconv.Atoi(contractIDStr)
	sig, err := h.Service.GetByContractID(contractID)
	if err != nil {
		http.Error(w, "not found", http.StatusOK)
		return
	}
	json.NewEncoder(w).Encode(sig)
}

func (h *SignatureHandler) GetContractsByCompanyID(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get(":id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}

	sigs, err := h.Service.GetContractsByCompanyID(companyID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(sigs)
}

func (h *SignatureHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SignatureHandler) Sign(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ContractID  int    `json:"contract_id"`
		ClientName  string `json:"client_name"`
		ClientIIN   string `json:"client_iin"`
		ClientPhone string `json:"client_phone"`
		Method      string `json:"method"`
		CompanyID   int    `json:"company_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := h.Service.Sign(input.ContractID, input.ClientName, input.ClientIIN, input.ClientPhone, input.Method, input.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SignatureHandler) ServeSignedPDFByID(w http.ResponseWriter, r *http.Request) {
	// Получаем ID подписи из query
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Получаем подпись из БД
	signature, err := h.Service.GetByID(id)
	if err != nil || signature == nil {
		http.Error(w, "signature not found", http.StatusNotFound)
		return
	}

	fmt.Println("Signature:", signature)

	// Проверяем путь к файлу
	filePath := filepath.FromSlash(signature.SignFilePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("File does not exist at path:", filePath)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// Заголовки для отображения PDF в браузере
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(filePath)+"\"")
	http.ServeFile(w, r, filePath)
}
