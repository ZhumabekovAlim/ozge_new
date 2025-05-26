package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
)

type SignatureFieldValueHandler struct {
	Service *services.SignatureFieldValueService
}

func NewSignatureFieldValueHandler(service *services.SignatureFieldValueService) *SignatureFieldValueHandler {
	return &SignatureFieldValueHandler{Service: service}
}

// POST /signature-fields/bulk
func (h *SignatureFieldValueHandler) CreateAll(w http.ResponseWriter, r *http.Request) {
	var input []models.SignatureFieldValue
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || len(input) == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for _, val := range input {
		if val.SignatureID == 0 || val.ContractFieldID == 0 {
			http.Error(w, "Missing fields", http.StatusBadRequest)
			return
		}
		if err := h.Service.Create(&val); err != nil {
			http.Error(w, "Failed to save field value", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /signature-fields/signature/:id
func (h *SignatureFieldValueHandler) GetBySignatureID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Error(w, "Invalid signature ID", http.StatusBadRequest)
		return
	}
	list, err := h.Service.GetBySignatureID(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(list)
}
