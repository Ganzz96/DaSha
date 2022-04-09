package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ganzz96/dasha/internal/agent_manager"
	"github.com/ganzz96/dasha/internal/agent_manager/config"
	"github.com/ganzz96/dasha/internal/agent_manager/storage"
	"github.com/ganzz96/dasha/internal/common/log"
	"github.com/ganzz96/dasha/internal/common/router"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	cli := cobra.Command{
		Use:   "agent-manager",
		Short: "DaSha Agent Manager",
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
		Short: "Run Agent Manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.WithStack(err)
			}

			return initAndRun(cfgPath)
		},
	}

	cmd.Flags().StringP("config", "c", "", "/path/to/config.yaml")
	_ = cmd.MarkFlagFilename("config")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func initAndRun(cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return errors.WithStack(err)
	}

	logger := log.New()
	router := router.NewRouter(logger)

	socket, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(cfg.AgentPusher.Host),
		Port: cfg.AgentPusher.Port,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer socket.Close()

	db, err := storage.New(logger, cfg.DBPath)
	if err != nil {
		return errors.WithStack(err)
	}

	gwPusher := agent_manager.NewAgentPusher(socket)
	go gwPusher.Serve()

	agentController := agent_manager.New(db, gwPusher)
	agentController.RegisterAPI(router)

	if err := http.ListenAndServe(cfg.HostPort, router); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
