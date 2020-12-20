package payment

import "payments-handler/entity"

type Provider struct {
	ID   int64  `json:"id"`
	Name string `json:"name,omitempty"`
}

type Option struct {
	Provider  Provider `json:"provider"`
	ButtonURL string   `json:"button_url"`
}

type AggregatorMock struct{}

func (a AggregatorMock) GetPaymentOptionsByProduct(_ entity.Product) ([]Option, error) {
	return []Option{
		{
			Provider: Provider{
				ID:   0,
				Name: "PayPal",
			},
			ButtonURL: "https://payments.paypal.com/button",
		},
		{
			Provider: Provider{
				ID:   1,
				Name: "Apple Pay",
			},
			ButtonURL: "https://payments.applepay.com/button",
		},
		{
			Provider: Provider{
				ID:   2,
				Name: "Google Pay",
			},
			ButtonURL: "https://payments.googlepay.com/button",
		},
		{
			Provider: Provider{
				ID:   3,
				Name: "Stripe",
			},
			ButtonURL: "https://payments.stripe.com/button",
		},
	}, nil
}
