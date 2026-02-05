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

// =====================================================
// /product
// GET    /product
// GET    /product?name=indomie&active=true
// POST   /product
// =====================================================
func (h *ProductHandler) Products(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// -----------------------------
	// GET /product (with filter)
	// -----------------------------
	case http.MethodGet:
		name := r.URL.Query().Get("name")

		var active *bool
		if v := r.URL.Query().Get("active"); v != "" {
			b, err := strconv.ParseBool(v)
			if err != nil {
				http.Error(w, "invalid active value (true/false)", http.StatusBadRequest)
				return
			}
			active = &b
		}

		// 🔍 FILTER MODE
		if name != "" || active != nil {
			products, err := h.service.Search(r.Context(), name, active)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(products)
			return
		}

		// 🔁 NORMAL GET ALL
		products, err := h.service.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)

	// -----------------------------
	// POST /product
	// -----------------------------
	case http.MethodPost:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.service.Create(r.Context(), &p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// =====================================================
// /product/{id}
// GET    /product/{id}
// PUT    /product/{id}
// DELETE /product/{id}
// =====================================================
func (h *ProductHandler) ProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	switch r.Method {

	// -----------------------------
	// GET /product/{id}
	// -----------------------------
	case http.MethodGet:
		product, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(product)

	// -----------------------------
	// PUT /product/{id}
	// -----------------------------
	case http.MethodPut:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		p.ID = id

		if err := h.service.Update(r.Context(), &p); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(p)

	// -----------------------------
	// DELETE /product/{id}
	// -----------------------------
	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
