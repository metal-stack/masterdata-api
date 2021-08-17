package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	gw "github.com/metal-stack/masterdata-api/api/v1"
	"github.com/metal-stack/v"
)

const (
	moduleName = "masterdata-api-gateway"
)

var (
	logger *zap.SugaredLogger
)

var rootCmd = &cobra.Command{
	Use:     moduleName,
	Short:   "the masterdata-api gateway provides a REST API next to the gRPC API",
	Version: v.V.String(),
	RunE: func(cmd *cobra.Command, args []string) error {
		initLogging()
		return run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		if logger != nil {
			logger.Fatalw("failed executing root command", "error", err)
		} else {
			log.Fatalf("failed executing root command: %s", err)
		}
	}
}

func initConfig() {
	viper.SetEnvPrefix("MASTERDATA_API_GATEWAY")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().IntP("port", "", 9000, "the port to serve on")
	rootCmd.Flags().StringP("base-path", "", "/", "the base path of the server")
	rootCmd.Flags().StringP("masterdata-api-address", "", "masterdata-api:9000", "the address of the masterdata-api")

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		log.Fatalf("unable to construct root command: %s", err)
	}
}

func initLogging() {
	level := zap.InfoLevel

	if viper.IsSet("log-level") {
		err := level.UnmarshalText([]byte(viper.GetString("log-level")))
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	l, err := cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger = l.Sugar()
}

func run() error {
	basePath := viper.GetString("base-path")
	if !strings.HasPrefix(basePath, "/") || !strings.HasSuffix(basePath, "/") {
		logger.Fatal("base path must start and end with a slash")
	}

	port := viper.GetInt("port")
	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	registerFns := []func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error){
		gw.RegisterProjectServiceHandlerFromEndpoint,
		gw.RegisterTenantServiceHandlerFromEndpoint,
	}
	for _, registerFn := range registerFns {
		err := registerFn(context.Background(), gwmux, viper.GetString("masterdata-api-address"), opts)
		if err != nil {
			return err
		}
	}

	mux.Handle(basePath, gwmux)

	logger.Infow("starting masterdata-api-gateway", "version", v.V, "address", addr)

	return http.ListenAndServe(addr, mux)
}
