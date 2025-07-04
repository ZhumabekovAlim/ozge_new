package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	_ "github.com/bmizerany/pat"
	"net/http"
	"strconv"
	"strings"
)

type CompanyHandler struct {
	Service *services.CompanyService
}

type CheckPhoneRequest struct {
	Phone string `json:"phone"`
}

type CheckPhoneResponse struct {
	ID int `json:"id"`
}

func NewCompanyHandler(s *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{Service: s}
}

// POST /companies/register
func (h *CompanyHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.Company
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil ||
		input.Name == "" || input.Password == "" || input.IIN == "" || (input.Email == "" && input.Phone == "") {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	company, err := h.Service.Register(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

func (h *CompanyHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Phone    string `json:"phone"`    // email или phone
		Password string `json:"password"` // пароль
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil ||
		input.Phone == "" || input.Password == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	company, err := h.Service.Login(input.Phone, input.Password)
	if err != nil {
		http.Error(w, "invalid login or password", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

// GET /companies
func (h *CompanyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	opts := services.CompanyListOptions{}

	// cursor
	if cursorStr := query.Get("cursor"); cursorStr != "" {
		if id, err := strconv.Atoi(cursorStr); err == nil {
			opts.CursorID = id
		}
	}

	// limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			opts.Limit = l
		}
	}
	if opts.Limit == 0 {
		opts.Limit = 10
	}

	// direction
	direction := strings.ToLower(query.Get("direction"))
	if direction != "prev" {
		direction = "next"
	}
	opts.Direction = direction

	// filters
	if search := query.Get("search"); search != "" {
		opts.Search = search
	}
	if idStr := query.Get("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			opts.FilterID = &id
		}
	}
	if name := query.Get("name"); name != "" {
		opts.FilterName = name
	}
	if email := query.Get("email"); email != "" {
		opts.FilterEmail = email
	}

	// sorting
	opts.SortBy = query.Get("sort")
	opts.Order = query.Get("order")

	// вызов сервиса
	companies, err := h.Service.List(opts)
	if err != nil {
		http.Error(w, "cannot get companies", http.StatusInternalServerError)
		return
	}

	var nextCursor, prevCursor int
	if len(companies) > 0 {
		prevCursor = companies[0].ID
		nextCursor = companies[len(companies)-1].ID
	}

	// ответ
	response := map[string]interface{}{
		"data":        companies,
		"next_cursor": nextCursor,
		"prev_cursor": prevCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func (h *CompanyHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
	var req CheckPhoneRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Phone == "" {
		http.Error(w, "phone required", http.StatusBadRequest)
		return
	}

	companyID, err := h.Service.GetCompanyIDByPhone(req.Phone)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if companyID == 0 {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}

	resp := CheckPhoneResponse{ID: companyID}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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

func (h *CompanyHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.OldPassword == "" || input.NewPassword == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Service.ChangePassword(id, input.OldPassword, input.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CompanyHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid company ID", http.StatusBadRequest)
		return
	}

	var input struct {
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.NewPassword == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Service.ResetPassword(id, input.NewPassword); err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
