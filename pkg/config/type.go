package config

import "time"

// Config struct to store all config
type (
	Config struct {
		HostList    []string      `json:"host_list"`
		PingTimeout time.Duration `json:"ping_timeout"`
		CronConfig  CronConfig    `json:"cron_config"`
	}

	CronConfig struct {
		HealthCheckAll string `json:"health_check_all"`
	}
)
