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

// GetDashboardSummary handles GET /dashboard/summary and returns aggregated
// statistics in JSON format.
func (h *StatisticsHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.Service.GetDashboardSummary()
	if err != nil {
		http.Error(w, "cannot get dashboard summary", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(summary)
}
