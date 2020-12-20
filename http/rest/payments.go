package rest

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"payments-handler/entity"
	"payments-handler/service/payment"
	"strings"
)

type PaymentsModule struct {
	Store      Store
	Aggregator PaymentProvidersAggregator
}

type Store interface {
	GetProductByID(string) (*entity.Product, error)
}

type PaymentProvidersAggregator interface {
	GetPaymentOptionsByProduct(entity.Product) ([]payment.Option, error)
}

func (m PaymentsModule) RegisterRoutes(r *mux.Router) {
	r = r.NewRoute().PathPrefix("/payments").Subrouter()

	r.HandleFunc("/options", m.GetOptionsByProductID).
		Methods(http.MethodGet)
}

type PaymentOptionsResponse struct {
	Product entity.Product   `json:"product"`
	Options []payment.Option `json:"options,omitempty"`
}

func (m PaymentsModule) GetOptionsByProductID(w http.ResponseWriter, r *http.Request) {
	var productID = strings.TrimSpace(r.FormValue("product_id"))

	if productID == "" {
		respondJSON(w, http.StatusBadRequest, failResponseBody{"product ID required"})
		return
	}

	product, err := m.Store.GetProductByID(productID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve product")
		respondJSON(w, http.StatusInternalServerError, failResponseBody{"something went wrong"})
		return
	}

	if product == nil {
		respondJSON(w, http.StatusNotFound, failResponseBody{"product does not exist"})
		return
	}

	options, err := m.Aggregator.GetPaymentOptionsByProduct(*product)
	if err != nil {
		log.Error().Err(err).Msg("Failed to aggregate payment options")
		respondJSON(w, http.StatusServiceUnavailable, failResponseBody{"third-party error happened"})
		return
	}

	respondJSON(w, http.StatusOK, PaymentOptionsResponse{
		Product: *product,
		Options: options,
	})
}
