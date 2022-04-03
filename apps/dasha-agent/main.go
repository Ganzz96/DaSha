package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ganzz96/dasha-agent/internal/agent"
	"github.com/ganzz96/dasha-agent/internal/clients"
	"github.com/ganzz96/dasha-agent/internal/config"
	"github.com/ganzz96/dasha-agent/internal/log"
)

func main() {
	cli := cobra.Command{
		Use:   "dasha-agent",
		Short: "DaSha - Data Sharing Decentralized Cloud Storage Agent",
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
		Short: "Run DaSha Agent",
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

func initAndRun(cfgPath string) error {
	cfgController := config.New()
	logger := log.New()

	cfg, err := cfgController.Load(cfgPath)
	if err != nil {
		return errors.WithStack(err)
	}

	meta, err := cfgController.LoadMeta()
	if err != nil {
		return errors.WithStack(err)
	}

	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer socket.Close()

	dashaClient, err := clients.NewDashaManagerClient(
		logger,
		socket,
		cfg.DashaManagerClient.HTTPHostPort,
		cfg.DashaManagerClient.UDPHostPort,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	agentController := agent.New(dashaClient, cfgController)
	if meta == nil {
		if err := agentController.Register(); err != nil {
			return errors.WithStack(err)
		}

		meta, err = cfgController.LoadMeta()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	reporter := agent.NewAgentReporter(logger, dashaClient, meta, time.Second*time.Duration(cfg.ReportInternavalInSec))
	go reporter.Up()

	dataMonitor := agent.NewDataMonitor(logger, socket)
	go dataMonitor.Up()

	ch := make(chan struct{})
	<-ch

	return nil
}
