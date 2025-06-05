package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	_ "github.com/bmizerany/pat"
	"net/http"
	"strconv"
)

type CompanyHandler struct {
	Service *services.CompanyService
}

func NewCompanyHandler(s *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{Service: s}
}

// POST /companies/register
func (h *CompanyHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.Company
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Name == "" || input.Email == "" || input.Phone == "" || input.Password == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	err = h.Service.Register(&input)
	if err != nil {
		http.Error(w, "could not create company", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// POST /companies/login
func (h *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Phone == "" || input.Password == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	company, err := h.Service.Login(input.Phone, input.Password)
	if err != nil {
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	company.Password = ""
	json.NewEncoder(w).Encode(company)
}

// GET /companies
func (h *CompanyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	companies, err := h.Service.List()
	if err != nil {
		http.Error(w, "cannot get companies", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(companies)
}

// GET /companies/id/:id
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}
	company, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(company)
}

// GET /companies/phone/:phone
func (h *CompanyHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get(":phone")
	if phone == "" {
		http.Error(w, "phone required", http.StatusBadRequest)
		return
	}
	company, err := h.Service.GetByPhone(phone)
	if err != nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(company)
}

// PUT /companies/:id
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}

	var input models.Company
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

// DELETE /companies/:id
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
