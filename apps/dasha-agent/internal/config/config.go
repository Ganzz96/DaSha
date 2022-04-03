package config

type AgentConfig struct {
	DashaManagerClient struct {
		UDPHostPort  string `json:"udp_hostport"`
		HTTPHostPort string `json:"http_hostport"`
	} `json:"dasha_manager_client"`
	ReportInternavalInSec int `json:"report_interval_in_sec"`
}

type AgentMeta struct {
	AgentID string `json:"agent_id"`
}
