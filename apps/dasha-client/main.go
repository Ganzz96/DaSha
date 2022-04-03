package main

import (
	"fmt"
	"os"

	"github.com/ganzz96/dasha-client/internal/clients"
	"github.com/ganzz96/dasha-client/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type UploadRequest struct {
	AgentID string
}

type UploadResponse struct {
	Conn string
}

func main() {
	cli := cobra.Command{
		Use:   "dasha-client",
		Short: "DaSha - Data Sharing Decentralized Cloud Storage",
	}

	cli.AddCommand(uploadCommand())

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func uploadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload file on DaSha",
		RunE: func(cmd *cobra.Command, args []string) error {
			agent, err := cmd.Flags().GetString("agent")
			if err != nil {
				return errors.WithStack(err)
			}

			config, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.WithStack(err)
			}

			message, err := cmd.Flags().GetString("message")
			if err != nil {
				return errors.WithStack(err)
			}

			return upload(agent, config, message)
		},
	}

	cmd.Flags().StringP("agent", "a", "", "agent-uuid-string")
	cmd.Flags().StringP("config", "c", "", "agent-uuid-string")
	cmd.Flags().StringP("message", "m", "", "any message string")

	_ = cmd.MarkFlagRequired("agent")
	_ = cmd.MarkFlagRequired("config")
	_ = cmd.MarkFlagRequired("message")

	return cmd
}

func upload(agentID string, cfgPath string, message string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return errors.WithStack(err)
	}

	managerClient, err := clients.NewDashaManagerClient(cfg.HostPort)
	if err != nil {
		return errors.WithStack(err)
	}

	agentInfo, err := managerClient.Upload(clients.UploadRequest{
		AgentID: agentID,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	agentClient, err := clients.NewAgentClient(agentInfo.Conn)
	if err != nil {
		return errors.WithStack(err)
	}

	return agentClient.Transmit(message)
}
