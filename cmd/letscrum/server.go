package main

import (
	"context"
	"github.com/daocloud/skoala/app/hive/internal/gateway"
	"github.com/daocloud/skoala/app/hive/internal/gateway/grpc"
	"github.com/daocloud/skoala/pkg/log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run both grpc and http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		log, err := log.New()
		if err != nil {
			return err
		}
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			log.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
		})

		stop := make(chan struct{})
		defer waitSignal(stop)

		if err := grpc.Run(context.Background(), "tcp", ":9090"); err != nil {
			log.Fatalf("grpc start error: %s", err)
		}
		opts := gateway.Options{
			Addr: ":8081",
			GRPCServer: gateway.Endpoint{
				Network: "tcp",
				Addr:    "localhost:9090",
			},
			OpenAPIDir: "api/hive/v1",
		}
		if err := gateway.Run(context.Background(), opts); err != nil {
			log.Fatalf("grpc gateway start error: %s", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func waitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}
