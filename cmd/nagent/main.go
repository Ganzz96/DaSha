package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ccding/go-stun/stun"
	"github.com/ganzz96/dasha/internal/nagent/agent"
	"github.com/ganzz96/dasha/internal/nagent/clients"
	"github.com/ganzz96/dasha/internal/nagent/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {
	cli := cobra.Command{
		Use:   "nagent",
		Short: "DaSha Node Agent",
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
		Short: "Run Node Agent",
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
	cfgController := config.New()

	cfg, err := cfgController.Load(cfgPath)
	if err != nil {
		return errors.WithStack(err)
	}

	meta, err := cfgController.LoadMeta()
	if err != nil {
		return errors.WithStack(err)
	}

	cl := stun.NewClient()
	cl.SetServerAddr(cfg.StunServer)

	_, externalAddr, err := cl.Discover()
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("Resolved stun external addr", externalAddr.String())

	udpSocket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer udpSocket.Close()

	tcpAgentManagerAddr, err := net.ResolveTCPAddr("tcp", cfg.AgentManager.TCPHostPort)
	if err != nil {
		return errors.WithStack(err)
	}

	tcpSocket, err := net.DialTCP("tcp", nil, tcpAgentManagerAddr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tcpSocket.Close()

	agentManagerClient, err := clients.NewAgentManagerClient(
		cfg.AgentManager.HTTPHostPort,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	agentController := agent.New(agentManagerClient, cfgController)
	if meta == nil {
		if err := agentController.Register(); err != nil {
			return errors.WithStack(err)
		}

		meta, err = cfgController.LoadMeta()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	reporter := agent.NewAgentReporter(agentManagerClient, meta, externalAddr.String(), time.Second*time.Duration(cfg.ReportInternavalInSec))
	go reporter.Up()

	agentGw := clients.NewAgentGateway(udpSocket)
	go agentGw.Serve()

	connManager := agent.NewConnManager(tcpSocket, agentGw)
	connManager.Serve()

	return nil
}
