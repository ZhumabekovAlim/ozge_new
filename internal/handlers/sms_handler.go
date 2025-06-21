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
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	code, err := h.Service.SendVerificationCode(req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"code":   code, // Для теста можно вернуть, но на проде скрывать!
	})
}
