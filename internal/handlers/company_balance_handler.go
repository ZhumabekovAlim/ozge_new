package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
)

type CompanyBalanceHandler struct {
	Service *services.CompanyBalanceService
}

func NewCompanyBalanceHandler(service *services.CompanyBalanceService) *CompanyBalanceHandler {
	return &CompanyBalanceHandler{Service: service}
}

func (h *CompanyBalanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cb models.CompanyBalance
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(&cb); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cb)
}

func (h *CompanyBalanceHandler) GetByCompanyID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	cb, err := h.Service.GetByCompanyID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cb)
}

func (h *CompanyBalanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)

	var cb models.CompanyBalance
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	cb.CompanyID = id

	if err := h.Service.Update(&cb); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CompanyBalanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
