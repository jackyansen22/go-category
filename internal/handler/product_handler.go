package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// /api/produk
func (h *ProductHandler) Products(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		data, err := h.service.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(data)

	case http.MethodPost:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid body", 400)
			return
		}

		if err := h.service.Create(r.Context(), &p); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "method not allowed", 405)
	}
}

// /api/produk/{id}
func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	switch r.Method {

	case http.MethodGet:
		p, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodPut:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid body", 400)
			return
		}
		p.ID = id

		if err := h.service.Update(r.Context(), &p); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", 405)
	}
}
