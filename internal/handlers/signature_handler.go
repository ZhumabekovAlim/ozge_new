package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	if err := h.Service.Create(&input); err != nil {
		http.Error(w, "failed to create signature", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *SignatureHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	sig, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(sig)
}

func (h *SignatureHandler) GetByContractID(w http.ResponseWriter, r *http.Request) {
	contractIDStr := r.URL.Query().Get(":id")
	contractID, _ := strconv.Atoi(contractIDStr)
	sig, err := h.Service.GetByContractID(contractID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(sig)
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

	err := h.Service.Sign(input.ContractID, input.ClientName, input.ClientIIN, input.ClientPhone, input.Method, input.CompanyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
