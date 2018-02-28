//go:generate gorunpkg github.com/vektah/gqlgen -out graph/generated.go -models graph/models.go -schema ./schema.gql -package graph -typemap types.json

package main

import (
	"net/http"

	"github.com/Teddy-Schmitz/graphqlCoin/rpcclient"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "net/http/pprof"

	"github.com/Teddy-Schmitz/graphqlCoin/graph"
	"github.com/gorilla/mux"
	"github.com/sethgrid/pester"
	"github.com/spf13/cobra"
	"github.com/vektah/gqlgen/handler"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "graphqlCoin",
		Short: "Graphql server for the blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			serverCmd()
		},
	}

	rootCmd.Execute()
}

func serverCmd() {
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if viper.GetBool("profiler") {
		go http.ListenAndServe("localhost:6060", nil)
	}

	p := pester.New()
	p.Concurrency = 1
	p.MaxRetries = 3
	p.Backoff = pester.ExponentialBackoff

	resolver := &graph.Resolver{
		Client: &rpcclient.Client{
			Client:   p,
			Host:     viper.GetString("daemon"),
			User:     viper.GetString("rpcuser"),
			Password: viper.GetString("rpcpassword"),
		},
	}

	r := mux.NewRouter()
	r.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(resolver)))

	logrus.Infoln("server started")
	logrus.Fatalln(http.ListenAndServe(":5050", r))
}

func init() {
	flag.StringP("daemon", "d", "", "Address and Port of coin daemon to use")
	flag.StringP("rpcuser", "u", "", "JSON-RPC User")
	flag.StringP("rpcpassword", "p", "", "JSON-RPC Password")

	flag.Bool("debug", false, "Enable debug logging")
	flag.Bool("profiler", false, "Enable pprof")

	viper.SetConfigName("graphqlcoin")
	viper.AddConfigPath("$HOME/.graphqlcoin")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.Fatalln(err)
		}
	}

	viper.AutomaticEnv()

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
}
