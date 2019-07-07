package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/toxiproxy/pkg/proxy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "toxiproxyd",
	Short: "A TCP proxy to simulate network and system conditions for chaos and resiliency testing",
	Run: func(cmd *cobra.Command, args []string) {
		server := proxy.NewServer()
		if config := viper.GetString("config"); config != "" {
			server.PopulateConfig(config)
		}

		// Handle SIGTERM to exit cleanly
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGTERM)
		go func() {
			<-signals
			os.Exit(0)
		}()

		server.Listen(viper.GetString("bind"), viper.GetString("port"))
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithField("error", err).Error("failed to run command")
		os.Exit(1)
	}
}

func addConfigBindings(cmd *cobra.Command) {
	cmd.Flags().String(
		"bind",
		"0.0.0.0",
		"address for toxiproxy to listen on",
	)
	viper.BindEnv("bind")
	viper.BindPFlag("bind", cmd.Flags().Lookup("bind"))

	cmd.Flags().String(
		"port",
		"8474",
		"port for toxiproxy to listen on",
	)
	viper.BindEnv("port")
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))

	cmd.Flags().String(
		"config",
		"",
		"JSON file containing proxies to create on startup",
	)
	viper.BindEnv("config")
	viper.BindPFlag("config", cmd.Flags().Lookup("config"))

	cmd.Flags().Int64(
		"seed",
		time.Now().UTC().UnixNano(),
		"seed for randomizing toxics with",
	)
	viper.BindEnv("seed")
	viper.BindPFlag("seed", cmd.Flags().Lookup("seed"))
}

func init() {
	viper.SetEnvPrefix("toxiproxy")
	cobra.OnInitialize(viper.AutomaticEnv)
	addConfigBindings(rootCmd)

	rand.Seed(viper.GetInt64("seed"))
}
