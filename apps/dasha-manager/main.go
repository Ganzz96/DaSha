package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ganzz96/dasha-manager/internal/agents"
	"github.com/ganzz96/dasha-manager/internal/agents/monitor"
	"github.com/ganzz96/dasha-manager/internal/config"
	"github.com/ganzz96/dasha-manager/internal/filestream"
	"github.com/ganzz96/dasha-manager/internal/log"
	"github.com/ganzz96/dasha-manager/internal/router"
	"github.com/ganzz96/dasha-manager/internal/storage"
)

func main() {
	cli := cobra.Command{
		Use:   "dasha-manager",
		Short: "DaSha - Data Sharing Decentralized Cloud Storage",
	}

	cli.AddCommand(runCommand())

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run instance of DaSha manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.WithStack(err)
			}

			return initAndRun(configPath)
		},
	}

	cmd.Flags().StringP("config", "c", "", "/path/to/config.yaml")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func initAndRun(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return errors.WithStack(err)
	}

	logger := log.New()
	router := router.NewRouter(logger)

	db, err := storage.New(logger, cfg.DBPath)
	if err != nil {
		return errors.WithStack(err)
	}

	agentController := agents.New(db)
	agentController.RegisterAPI(router)

	filestreamController := filestream.New(agentController)
	filestreamController.RegisterAPI(router)

	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(cfg.AgentMonitor.Host),
		Port: cfg.AgentMonitor.Port,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer socket.Close()

	agentMonitor := monitor.New(logger, socket, agentController)
	go agentMonitor.Serve()

	if err := http.ListenAndServe(cfg.HostPort, router); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
