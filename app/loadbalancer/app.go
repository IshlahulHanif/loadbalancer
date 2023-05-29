package main

import (
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/api/httpapi"
	cronjob "github.com/loadbalancer/job/cron"
	"github.com/loadbalancer/pkg/config"
	"net/http"
	"time"
)

func main() {
	var (
		err error
	)

	defer func() {
		if err != nil {
			fmt.Println(poneglyph.GetErrorLogMessage(err))
		}
	}()

	// init config
	conf := config.Config{
		HostList: []string{ //TODO: read from config file
			"http://127.0.0.1:8081",
			"http://127.0.0.1:8082",
			"http://127.0.0.1:8083",
		},
		PingTimeout:    2 * time.Second,
		RequestTimeout: 2 * time.Second,
		CronConfig: config.CronConfig{
			HealthCheckAll: "@every 10s",
		},
	}

	// init poneglyph settings
	poneglyph.SetProjectName("loadbalancer")
	poneglyph.SetIsPrintFromContentRoot(true)
	poneglyph.SetIsPrintFunctionName(true)
	poneglyph.SetIsPrintNewline(true)
	poneglyph.SetIsUseTabSeparator(false)

	// TODO: tidy up

	// init http api
	httpApi, err := httpapi.GetInstance(conf)
	if err != nil {
		err = poneglyph.Trace(err)
		return
	}

	// init cron job scheduler
	cronJob, err := cronjob.GetInstance(conf)
	if err != nil {
		err = poneglyph.Trace(err)
		return
	}

	// Create two routers for different ports
	requestRouter := http.NewServeMux()

	hostManagementRouter := http.NewServeMux()

	// register handlers for request forwarder
	requestRouter.HandleFunc("/", httpApi.HandlerForwardRequest)

	// register handlers for host management
	hostManagementRouter.HandleFunc("/host/manage", httpApi.HandlerManageHostRequest)

	var errChan = make(chan error)

	// Start the requestRouter server on port 8080
	go func() {
		errServer := http.ListenAndServe(":8080", requestRouter)
		if errServer != nil {
			errChan <- poneglyph.Trace(errServer)
		}
	}()

	// Start the hostManagementRouter server on port 9090
	go func() {
		errServer := http.ListenAndServe(":9090", hostManagementRouter)
		if errServer != nil {
			errChan <- poneglyph.Trace(errServer)
		}
	}()

	// Start cron scheduler
	err = cronJob.StartAllCronJobScheduler()
	if err != nil {
		err = poneglyph.Trace(err)
		return
	}

	err = <-errChan
	if err != nil {
		err = poneglyph.Trace(err)
		// TODO: consume all errChan to avoid memory leak
	}

	// TODO: make sure to Shutdown the servers gracefully
}
