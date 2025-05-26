package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"OzgeContract/internal/services"
)

type StatisticsHandler struct {
	Service *services.StatisticsService
}

func NewStatisticsHandler(service *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{Service: service}
}

func (h *StatisticsHandler) GetCompanyStats(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	companyID, _ := strconv.Atoi(idStr)

	stats, err := h.Service.GetCompanyStats(companyID)
	if err != nil {
		http.Error(w, "cannot get stats", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
