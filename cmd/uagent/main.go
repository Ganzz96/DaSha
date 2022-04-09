package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ccding/go-stun/stun"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ganzz96/dasha/internal/common/log"
	"github.com/ganzz96/dasha/internal/common/router"
	"github.com/ganzz96/dasha/internal/uagent"
	"github.com/ganzz96/dasha/internal/uagent/clients"
	"github.com/ganzz96/dasha/internal/uagent/config"
)

func main() {
	cli := cobra.Command{
		Use:   "uagent",
		Short: "DaSha User Agent",
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
		Short: "Run User Agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.WithStack(err)
			}

			return initAndRun(cfgPath)
		},
	}

	cmd.Flags().StringP("config", "c", "", "agent-uuid-string")
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

	socket, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0")})
	if err != nil {
		return errors.WithStack(err)
	}
	defer socket.Close()

	cl := stun.NewClient()
	cl.SetServerAddr(cfg.StunServer)

	nat, externalAddr, err := cl.Discover()
	if err != nil {
		return errors.WithStack(err)
	}

	amClient, err := clients.NewAgentManagerClient(cfg.AgentManagerEndpoint)
	if err != nil {
		return errors.WithStack(err)
	}

	agentGateway := clients.NewAgentGateway(socket)
	go agentGateway.Serve()

	uagent, err := uagent.New(nat, externalAddr, amClient, agentGateway)
	if err != nil {
		return errors.WithStack(err)
	}

	uagent.RegisterAPI(router)

	return http.ListenAndServe(cfg.HostPort, router)
}
