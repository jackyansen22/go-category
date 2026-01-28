package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/service"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// /categories
func (h *CategoryHandler) Categories(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 CATEGORIES HANDLER HIT:", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		categories, err := h.service.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(categories)

	case http.MethodPost:
		var c model.Category
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.service.Create(r.Context(), &c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("DEBUG CATEGORY:", c)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c) // ❗ HARUS c

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// /categories/{id}
func (h *CategoryHandler) CategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		c, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(c)

	case http.MethodPut:
		var c model.Category
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		c.ID = id

		if err := h.service.Update(r.Context(), &c); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(c)

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
