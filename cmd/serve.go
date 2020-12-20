package cmd

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"payments-handler/http/rest"
	"payments-handler/service/payment"
	"payments-handler/service/store"
)

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Serve HTTP requests",
	RunE:  runServe,
}

const (
	addressFlag  = "address"
	addressFlagL = "a"
)

var (
	address string
)

func init() {
	serve.Flags().
		StringVarP(&address, addressFlag, addressFlagL, "localhost:8080", "port to listen on")

	root.AddCommand(serve)
}

func runServe(_ *cobra.Command, _ []string) error {
	var paymentsModule = rest.PaymentsModule{
		Store:      store.Store{},
		Aggregator: payment.Aggregator{},
	}

	var router = mux.NewRouter()

	router.Use(
		rest.HandlePanicMiddleware,
		rest.DeferredResponseMiddleware,
		rest.ResponseTimeMiddleware,
		rest.GetServerNameMiddleware(hostname),
	)

	paymentsModule.RegisterRoutes(router)

	log.Info().
		Str("address", address).
		Msg("Listening and serving HTTP requests...")

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal().Err(err).Msg("Listen and serve failed")
	}

	return nil
}
