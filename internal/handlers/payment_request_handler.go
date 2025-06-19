package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type PaymentRequestHandler struct {
	Service *services.PaymentRequestService
}

func NewPaymentRequestHandler(service *services.PaymentRequestService) *PaymentRequestHandler {
	return &PaymentRequestHandler{Service: service}
}

func (h *PaymentRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Service.Create(&input); err != nil {
		http.Error(w, "create failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *PaymentRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	p, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
	w.WriteHeader(http.StatusOK)
}

func (h *PaymentRequestHandler) GetByCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	list, err := h.Service.GetByCompany(id)
	if err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list) // статус 200 по умолчанию
}

func (h *PaymentRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	cursorID := int(^uint(0) >> 1) // max int
	limit := 20

	if cursorStr := query.Get("cursor"); cursorStr != "" {
		if c, err := strconv.Atoi(cursorStr); err == nil {
			cursorID = c
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	list, err := h.Service.GetAll(r.Context(), cursorID, limit)
	if err != nil {
		log.Printf("handler error: %v", err)
		http.Error(w, "fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var nextCursor int
	if len(list) > 0 {
		nextCursor = list[len(list)-1].ID
	}

	log.Printf("returning %d payment requests", len(list))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        list,
		"next_cursor": nextCursor,
	})
}

func (h *PaymentRequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	var input models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	input.ID = id
	if err := h.Service.Update(&input); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *PaymentRequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
