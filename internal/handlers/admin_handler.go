package handlers

import (
	"OzgeContract/internal/models"
	"OzgeContract/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	Service *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{Service: service}
}

func (h *AdminHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.Admin
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Email == "" || input.Password == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Service.Register(&input); err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Email == "" || input.Password == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	admin, err := h.Service.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, "invalid login or password", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.List()
	if err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *AdminHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	a, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(a)
}

func (h *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	var input models.Admin
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

func (h *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Service.Delete(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, _ := strconv.Atoi(idStr)
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
