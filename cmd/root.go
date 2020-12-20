package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var hostname = "none"

func init() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get hostname")
	}

	hostname = host
}

var root = &cobra.Command{
	Use:   os.Args[0],
	Short: "FooBar API",
	Long:  "The cake is a lie",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// configuring global logger

		parsedLevel, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			return fmt.Errorf("failed to parse %s (%s): %w", logFlag, logLevel, err)
		}

		zerolog.SetGlobalLevel(parsedLevel)

		log.Info().Msgf("Log level set to %s", parsedLevel)

		return nil
	},
}

const (
	logFlag  = "loglevel"
	logFlagL = "l"
)

var (
	logLevel string
)

func init() {
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().Caller().Logger()

	root.PersistentFlags().
		StringVarP(&logLevel, logFlag, logFlagL, zerolog.DebugLevel.String(), "level of logging")
}

func Execute() {
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msgf("Failed to execute. Run '%s --help' to see usage", os.Args[0])
	}
}
