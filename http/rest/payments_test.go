package rest

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"payments-handler/entity"
	"payments-handler/service/payment"
	"payments-handler/service/store"
	"testing"
)

func TestGetOptionsByProductIDHandler(t *testing.T) {

	t.Run("NoProduct", func(t *testing.T) {
		var (
			req = getReqNoBody(http.MethodGet, "/payments/options")
			rec = httptest.NewRecorder()
		)

		h := PaymentsModule{
			Store:      &store.Store{},
			Aggregator: &payment.Aggregator{},
		}

		h.GetOptionsByProductID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InvalidProduct", func(t *testing.T) {
		var (
			req = getReqNoBody(http.MethodGet, "/payments/options?product_id=1")
			rec = httptest.NewRecorder()

			s = &storeMock{}
		)

		s.On("GetProductByID", mock.Anything).
			Return((*entity.Product)(nil), nil)

		h := PaymentsModule{
			Store:      s,
			Aggregator: &payment.Aggregator{},
		}

		h.GetOptionsByProductID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("StoreError", func(t *testing.T) {
		var (
			req = getReqNoBody(http.MethodGet, "/payments/options?product_id=5fdf8dd752da22ffcb1cf412")
			rec = httptest.NewRecorder()

			s = &storeMock{}
		)

		s.On("GetProductByID", mock.Anything).
			Return((*entity.Product)(nil), errors.New("failed to get product"))

		h := PaymentsModule{
			Store:      s,
			Aggregator: &payment.Aggregator{},
		}

		h.GetOptionsByProductID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AggregatorError", func(t *testing.T) {
		var (
			req = getReqNoBody(http.MethodGet, "/payments/options?product_id=5fdf8dd752da22ffcb1cf412")
			rec = httptest.NewRecorder()

			a = &paymentsAggregatorMock{}
		)

		a.On("GetPaymentOptionsByProduct", mock.Anything).
			Return(([]payment.Option)(nil), errors.New("failed to get options"))

		h := PaymentsModule{
			Store:      store.Store{},
			Aggregator: a,
		}

		h.GetOptionsByProductID(rec, req)

		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})

	t.Run("ValidResponse", func(t *testing.T) {
		var (
			req = getReqNoBody(http.MethodGet, "/payments/options?product_id=5fdf8dd752da22ffcb1cf412")
			rec = httptest.NewRecorder()

			s = &storeMock{}
			a = &paymentsAggregatorMock{}

			product = entity.Product{ID: "5fdf8dd752da22ffcb1cf412"}
			option  = payment.Option{
				Provider: payment.Provider{
					ID:   0,
					Name: "PayPal",
				},
				ButtonURL: "https://payments.paypal.com/button",
			}
		)

		s.On("GetProductByID", "5fdf8dd752da22ffcb1cf412").
			Return(&product, nil)

		a.On("GetPaymentOptionsByProduct", product).
			Return([]payment.Option{option}, nil)

		h := PaymentsModule{
			Store:      s,
			Aggregator: a,
		}

		h.GetOptionsByProductID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		const respBody = "{\"product\":{\"id\":\"5fdf8dd752da22ffcb1cf412\"},\"options\":[{\"provider\":{\"id\":0,\"name\":\"PayPal\"},\"button_url\":\"https://payments.paypal.com/button\"}]}\n"

		assert.Equal(t, respBody, string(rec.Body.Bytes()))
	})
}

func getReqNoBody(method, uri string) (r *http.Request) {
	r, err := http.NewRequest(method, uri, nil)
	if err != nil {
		panic(err)
	}
	return
}

type storeMock struct {
	mock.Mock
}

func (s *storeMock) GetProductByID(id string) (*entity.Product, error) {
	var args = s.Called(id)

	return args.Get(0).(*entity.Product), args.Error(1)
}

type paymentsAggregatorMock struct {
	mock.Mock
}

func (a *paymentsAggregatorMock) GetPaymentOptionsByProduct(p entity.Product) ([]payment.Option, error) {
	var args = a.Called(p)

	return args.Get(0).([]payment.Option), args.Error(1)
}
