package cmd

import (
	"context"
	"fmt"
	"github.com/letscrum/letscrum/internal/gateway"
	"github.com/letscrum/letscrum/internal/gateway/grpc"
	"github.com/letscrum/letscrum/pkg/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// cfgPath is the path to the EnvoyGateway configuration file.
	cfgPath string
)

func GetServerCommand() *cobra.Command {
	cmd := &cobra.Command{
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
			defer WaitSignal(stop)

			if err := grpc.Run(context.Background(), "tcp", ":9090"); err != nil {
				log.Fatalf("grpc start error: %s", err)
			}
			opts := gateway.Options{
				Addr: ":8081",
				GRPCServer: gateway.Endpoint{
					Network: "tcp",
					Addr:    "localhost:9090",
				},
				OpenAPIDir: "api/v1",
			}
			if err := gateway.Run(context.Background(), opts); err != nil {
				log.Fatalf("grpc gateway start error: %s", err)
			}

			return nil
		},
	}
	cobra.OnInitialize(InitConfig)
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./config/config.yaml", "config file (default is $HOME/.letscrum.yaml)")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cmd.PersistentFlags().StringVarP(&cfgPath, "config-path", "c", "",
		"The path to the configuration file.")

	return cmd
}

func WaitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".letscrum")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
