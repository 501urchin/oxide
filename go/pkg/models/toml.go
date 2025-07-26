package models

type Config struct {
	App          appConfig
	Server       serverConfig
	Firewall     firewallConfig     `toml:"firewall"`
	ReverseProxy reverseProxyConfig `toml:"reverse_proxy"`
	Healthcheck  healthCheckConfig  `toml:"healthcheck"`
	Secrets      secretsConfig      `toml:"secrets"`
}

type appConfig struct {
	Name                 string   `toml:"name"`
	SourceDirectory      string   `toml:"source_directory"`
	DeployType           string   `toml:"deploy_type"`
	BuildCommand         string   `toml:"build_command"`
	KeepPreviousReleases int      `toml:"keep_previous_releases"`
	BackupPrevious       bool     `toml:"backup_previous_release"`
	PreDeployScript      string   `toml:"pre_deploy_script"`
	PostDeployScript     string   `toml:"post_deploy_script"`
	Ports                []string `toml:"ports"`
}

type serverConfig struct {
	OS              string `toml:"os"`
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	User            string `toml:"user"`
	UsePassword     bool   `toml:"use_password"`
	Password        string `toml:"password"`
	SSHKeyInline    string `toml:"ssh_key_inline"`
	SSHKeyFile      string `toml:"ssh_key_file"`
	KnownHostsFile  string `toml:"known_hosts_file"`
	DeployDirectory string `toml:"deploy_directory"`
}

type firewallConfig struct {
	AllowedInboundPorts []int  `toml:"allowed_inbound_ports"`
	AllowedOutbound     string `toml:"allowed_outbound"`
	BlockICMPPing       bool   `toml:"block_icmp_ping"`
	LimitSSHRate        bool   `toml:"limit_ssh_rate"`
	AutoConfigure       bool   `toml:"auto_configure"`
}

type reverseProxyConfig struct {
	Enabled             bool   `toml:"enabled"`
	Type                string `toml:"type"`
	Domain              string `toml:"domain"`
	EnableHTTPS         bool   `toml:"enable_https"`
	RedirectHTTPToHTTPS bool   `toml:"redirect_http_to_https"`
}

type healthCheckConfig struct {
	Path           string `toml:"path"`
	ExpectedStatus int    `toml:"expected_status"`
}

type secretsConfig struct {
	LocalEnvFile  string `toml:"local_env_file"`
	RemoteEnvFile string `toml:"remote_env_file"`
}
