package cart

import (
	"fmt"
	"net/http"

	"github.com/KristianKjerstad/go-e-commerce-api/service/auth"
	"github.com/KristianKjerstad/go-e-commerce-api/types"
	"github.com/KristianKjerstad/go-e-commerce-api/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.HandleCheckout, h.userStore)).Methods(http.MethodPost)

}

func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	var cart types.CartCheckoutPayload
	err := utils.ParseJSON(r, &cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload %v", err))
		return
	}
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	products, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})

}
