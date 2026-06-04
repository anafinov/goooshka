package http

import (
	"encoding/json"
	"net/http"

	"student-cert-service/internal/domain"
	"student-cert-service/internal/service"
)

type Handlers struct {
	authService     *service.AuthService
	requestsService *service.RequestsService
}

func NewHandlers(authService *service.AuthService, requestsService *service.RequestsService) *Handlers {
	return &Handlers{
		authService:     authService,
		requestsService: requestsService,
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

type CreateCertRequest struct {
	Type string `json:"type"`
}

func (h *Handlers) CreateRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	certReq, err := h.requestsService.CreateRequest(userID, req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(certReq)
}

func (h *Handlers) GetRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	reqs, err := h.requestsService.GetMyRequests(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reqs == nil {
		reqs = []domain.Request{}
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	json.NewEncoder(w).Encode(reqs)
}
