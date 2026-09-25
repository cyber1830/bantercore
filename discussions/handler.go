package discussions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sahilverma/muze-go-backend/auth"
)

type Handler struct {
	s      *Service
	tokens *auth.TokenManager
}

func NewHandler(s *Service, t *auth.TokenManager) *Handler {
	return &Handler{s: s, tokens: t}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/health":
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/session":
		h.session(w)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/discussions":
		h.list(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/discussions":
		h.create(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) session(w http.ResponseWriter) {
	token, err := h.tokens.Issue("demo-user")
	if err != nil {
		http.Error(w, `{"error":"could not create session"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token, "userId": "demo-user"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	rows, err := h.s.List(r.Context(), limit)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rows)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	authorization := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := h.tokens.Verify(authorization)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var input struct {
		Topic string `json:"topic"`
		Body  string `json:"body"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	discussion, err := h.s.Create(r.Context(), user, input.Topic, input.Body)
	if err != nil {
		http.Error(w, `{"error":"validation error"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(discussion)
}
