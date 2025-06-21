package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type SignatureHandler struct {
	Service           *services.SignatureService
	FieldValueService *services.SignatureFieldValueService
}

func NewSignatureHandler(service *services.SignatureService, fvService *services.SignatureFieldValueService) *SignatureHandler {
	return &SignatureHandler{Service: service, FieldValueService: fvService}
}

func (h *SignatureHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	contractID, _ := strconv.Atoi(r.FormValue("contract_id"))
	var valueDTOs []models.SignatureFieldValueDTO
	if val := r.FormValue("field_values"); val != "" {
		_ = json.Unmarshal([]byte(val), &valueDTOs)
	}
	input := models.Signature{
		ContractID:  contractID,
		ClientName:  r.FormValue("client_name"),
		ClientIIN:   r.FormValue("client_iin"),
		ClientPhone: r.FormValue("client_phone"),
		Method:      r.FormValue("method"),
	}

	newID, err := h.Service.Create(&input)
	if err != nil {
		http.Error(w, "failed to create signature", http.StatusInternalServerError)
		return
	}
	input.ID = newID

	// Save additional field values if provided
	for _, v := range valueDTOs {
		fv := models.SignatureFieldValue{
			SignatureID:     newID,
			ContractFieldID: v.ContractFieldID,
			FieldValue:      v.FieldValue,
		}
		_ = h.FieldValueService.Create(&fv)
	}

	contract, err := h.Service.GetContractByID(input.ContractID)
	if err != nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}

	signDir := fmt.Sprintf("C:\\Users\\alimz\\GolandProjects\\OzgeContract\\uploads\\signatures\\company_%d", contract.CompanyID)
	if err := os.MkdirAll(signDir, 0755); err != nil {
		http.Error(w, "cannot create directory", http.StatusInternalServerError)
		return
	}

	signedPDF := filepath.Join(signDir, fmt.Sprintf("signed_final_%d.pdf", input.ID))
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	out, err := os.Create(signedPDF)
	if err != nil {
		http.Error(w, "cannot save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
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

func (h *SignatureHandler) GetSignaturesAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	opts := services.SignatureListOptions{}

	if cursorStr := query.Get("cursor"); cursorStr != "" {
		if id, err := strconv.Atoi(cursorStr); err == nil {
			opts.CursorID = id
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			opts.Limit = l
		}
	}
	if opts.Limit == 0 {
		opts.Limit = 20
	}

	if search := query.Get("search"); search != "" {
		opts.Search = search
	}

	if statusStr := query.Get("status"); statusStr != "" {
		if st, err := strconv.Atoi(statusStr); err == nil {
			opts.Status = &st
		}
	}

	if method := query.Get("method"); method != "" {
		opts.Method = method
	}

	opts.SortBy = query.Get("sort")
	opts.Order = query.Get("order")

	direction := query.Get("direction")
	if direction != "prev" {
		direction = "next"
	}
	opts.Direction = direction

	sigs, err := h.Service.GetSignaturesAll(opts)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var nextCursor, prevCursor int
	if len(sigs) > 0 {
		prevCursor = sigs[0].ID
		nextCursor = sigs[len(sigs)-1].ID
	}

	response := map[string]interface{}{
		"data":        sigs,
		"next_cursor": nextCursor,
		"prev_cursor": prevCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
