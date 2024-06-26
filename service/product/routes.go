package product

import (
	"fmt"
	"net/http"

	"github.com/KristianKjerstad/go-e-commerce-api/types"
	"github.com/KristianKjerstad/go-e-commerce-api/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodGet)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not get products"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)

}
