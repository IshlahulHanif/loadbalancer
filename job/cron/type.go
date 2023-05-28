package cronjob

import (
	"github.com/loadbalancer/pkg/config"
	"gopkg.in/robfig/cron.v2"
)

type Module struct {
	config     config.Config
	cronModule *cron.Cron //TODO: use interface for the cron module to enable mock
	service    service
}

type service struct {
	hostpool hostpoolMethod
}
