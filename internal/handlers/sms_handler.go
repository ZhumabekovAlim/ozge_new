package handlers

import (
	"OzgeContract/internal/services"
	"encoding/json"
	"net/http"
)

// SMSHandler handles HTTP requests for sending SMS messages.
type SMSHandler struct {
	Service *services.SMSService
}

// NewSMSHandler creates a new SMSHandler.
func NewSMSHandler(service *services.SMSService) *SMSHandler {
	return &SMSHandler{Service: service}
}

// Send handles POST /sms/send requests.
func (h *SMSHandler) Send(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Text  string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" || req.Text == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.SendSMS(req.Phone, req.Text); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
