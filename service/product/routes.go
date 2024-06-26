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
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Could not get products"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)

}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product = new(types.Product)
	err := utils.ParseJSON(r, &product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Could not get product from request"))
		return
	}
	err = h.store.CreateProduct(*product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Could not create product"))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, product)

}
