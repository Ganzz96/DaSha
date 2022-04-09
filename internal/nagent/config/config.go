package config

type AgentConfig struct {
	AgentManager struct {
		TCPHostPort  string `json:"tcp_hostport"`
		HTTPHostPort string `json:"http_hostport"`
	} `json:"agent_manager"`
	ReportInternavalInSec int    `json:"report_interval_in_sec"`
	StunServer            string `json:"stun_server"`
}

type AgentMeta struct {
	AgentID string `json:"agent_id"`
}
