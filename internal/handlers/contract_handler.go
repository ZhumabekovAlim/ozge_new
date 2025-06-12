package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
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
	if err := h.Service.Create(&input); err != nil {
		http.Error(w, "failed to create contract", http.StatusInternalServerError)
		return
	}
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

	// собираем contract
	contract := models.Contract{
		CompanyID:        input.CompanyID,
		TemplateID:       input.TemplateID,
		GeneratedPDFPath: input.GeneratedPDFPath,
		ClientFilled:     input.ClientFilled,
		Method:           input.Method,
		ContractToken:    uuid.New().String(),
	}

	// собираем contract fields (без contract_id)
	fields := make([]models.ContractField, 0)
	for _, f := range input.Fields {
		fields = append(fields, models.ContractField{
			FieldName: f.FieldName,
			FieldType: f.FieldType,
		})
	}

	// вызываем сервис
	err := h.Service.CreateWithFields(&contract, fields)
	if err != nil {
		http.Error(w, "failed to create contract and fields", http.StatusInternalServerError)
		return
	}
	// верни новый контракт + fields
	json.NewEncoder(w).Encode(contract)
}
