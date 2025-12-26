package orders

import (
	"log"
	"net/http"

	"github.com/Arnav-Soney/ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// reading payload
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call the service
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)
		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return to use the payloads/response with statusCode
	json.Write(w, http.StatusCreated, createdOrder)
}
