package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
)

type TariffPlanHandler struct {
	Service *services.TariffPlanService
}

func NewTariffPlanHandler(service *services.TariffPlanService) *TariffPlanHandler {
	return &TariffPlanHandler{Service: service}
}

func (h *TariffPlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tp models.TariffPlan
	if err := json.NewDecoder(r.Body).Decode(&tp); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(&tp); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tp)
}

func (h *TariffPlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *TariffPlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	tp, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(tp)
}

func (h *TariffPlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)

	var tp models.TariffPlan
	if err := json.NewDecoder(r.Body).Decode(&tp); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	tp.ID = id

	if err := h.Service.Update(&tp); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TariffPlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
