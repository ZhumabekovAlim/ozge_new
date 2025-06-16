package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	_ "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	_ "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
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

	// Сохраняем подпись
	newID, err := h.Service.Create(&input)
	if err != nil {
		http.Error(w, "failed to create signature", http.StatusInternalServerError)
		return
	}
	input.ID = newID

	// Получаем контракт по ID
	contract, err := h.Service.GetContractByID(input.ContractID)
	if err != nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}
	contractPDF := contract.GeneratedPDFPath

	signDir := fmt.Sprintf("uploads/signatures/company_%d", contract.CompanyID)
	if err := os.MkdirAll(signDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	signedPDF := filepath.Join(signDir, fmt.Sprintf("signed_final_%d.pdf", input.ID))
	err = AddSignatureToLastPage(
		contractPDF,       // исходный документ с контрактом
		signedPDF,         // финальный документ с подписью
		input.ClientName,  // ФИО
		input.ClientIIN,   // ИИН
		input.ClientPhone, // телефон
		"DejaVuSans.ttf",  // путь к ttf-шрифту с кириллицей!
	)
	if err != nil {
		http.Error(w, "signature page error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	input.SignFilePath = signedPDF
	_ = h.Service.UpdateSignFilePath(input.ID, signedPDF)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
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

func GenerateSignaturePDF(filename, name, iin, phone string, signedAt time.Time) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 18)
	pdf.AddPage()
	pdf.CellFormat(0, 20, "Подпись клиента", "", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("DejaVu", "", 14)
	pdf.SetXY(10, 40)
	pdf.MultiCell(0, 10, fmt.Sprintf(
		"ФИО физического лица: %s\nИНН: %s\nНомер телефона: %s\nПодписан: %s",
		name,
		iin,
		phone,
		signedAt.Format("2006-01-02 15:04:05"),
	), "", "L", false)
	return pdf.OutputFileAndClose(filename)
}

func AddSignatureToLastPage(
	inPDF, outPDF, fio, iin, phone, fontPath string,
) error {
	// 1. Узнаём количество страниц
	ctx, err := api.ReadContextFile(inPDF)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}
	pageCount := ctx.PageCount
	pageStr := strconv.Itoa(pageCount)

	// 2. Формируем текст для подписи
	text := fmt.Sprintf(
		"ФИО: %s\nИНН: %s\nТелефон: %s\nДата подписи: %s",
		fio, iin, phone, time.Now().Format("2006-01-02 15:04:05"),
	)

	// 3. Описание стиля
	desc := fmt.Sprintf("font:%s, points:13, pos:bl, offset:20 35, rot:0, op:1, fillc:0 0 0", fontPath)
	// Пример: font:/home/user/DejaVuSans.ttf,... или font:DejaVuSans.ttf,... если файл рядом с исполняемым

	// 4. Добавляем текст только на последнюю страницу
	return api.AddTextWatermarksFile(
		inPDF,
		outPDF,
		[]string{pageStr}, // только последняя страница
		false,             // не поверх, а “за” контентом (можешь поставить true если нужно наверх)
		text,
		desc,
		nil,
	)
}
