package cronjob

import (
	"context"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
)

func (m Module) handlerHealthCheckJob() {
	var (
		ctx = context.Background()
		err error
	)

	defer func() {
		if err != nil {
			fmt.Println(poneglyph.GetErrorLogMessage(err))
		}
	}()

	// execute health check
	res, err := m.service.hostpool.HealthCheckAllHost(ctx)
	if err != nil {
		err = poneglyph.Trace(err)
	}

	// print health check result
	// Print the health check result
	fmt.Println("Health Check Result:")
	fmt.Println("Healthy Hosts:")
	for _, host := range res.HealthyHosts {
		fmt.Printf("\t- %s\n", host)
	}

	fmt.Println("Down Hosts:")
	for _, host := range res.DownHosts {
		fmt.Printf("\t- %s\n", host)
	}
}
