package handlers

import (
	"OzgeContract/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

// SMSHandler handles HTTP requests for sending SMS/WhatsApp verification codes.
type SMSHandler struct {
	Service      *services.SMSService
	WhatsService *services.WhatsAppService
}

// NewSMSHandler creates a new SMSHandler.
func NewSMSHandler(sms *services.SMSService, wa *services.WhatsAppService) *SMSHandler {
	return &SMSHandler{Service: sms, WhatsService: wa}
}

// Send handles POST /sms/send requests.
func (h *SMSHandler) Send(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.Service == nil {
		http.Error(w, `{"error":"sms service not configured"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	code, err := h.Service.SendVerificationCode(req.Phone)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"code":   code, // Для теста — отправляем код на фронт
	})
}

// SendWhatsApp handles POST /whatsapp/send requests.
func (h *SMSHandler) SendWhatsApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.WhatsService == nil || h.WhatsService.SMSC == nil {
		http.Error(w, `{"error":"whatsapp service not configured"}`, http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Phone) == "" {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	code, err := h.WhatsService.SendVerificationCode(req.Phone)
	if err != nil {
		http.Error(w, `{"error":"`+strings.ReplaceAll(err.Error(), `"`, `'`)+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"code":   code, // оставить для тестов, в проде — убрать
	})
}
