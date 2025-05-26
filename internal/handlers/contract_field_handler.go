package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
)

type ContractFieldHandler struct {
	Service *services.ContractFieldService
}

func NewContractFieldHandler(service *services.ContractFieldService) *ContractFieldHandler {
	return &ContractFieldHandler{Service: service}
}

func (h *ContractFieldHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.ContractField
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(&input); err != nil {
		http.Error(w, "Failed to create field", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *ContractFieldHandler) GetByContractID(w http.ResponseWriter, r *http.Request) {
	contractID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}
	fields, err := h.Service.GetByContractID(contractID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(fields)
}
