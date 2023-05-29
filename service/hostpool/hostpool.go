package hostpool

import (
	"context"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	hostEntity "github.com/loadbalancer/entity/host"
)

func (s Service) ManageHost(ctx context.Context, req hostEntity.ManageHostReq) (err error) {
	switch req.Operation {
	case hostEntity.Add:
		for _, hostIP := range req.Data {
			err = s.usecase.hostpool.AddHost(ctx, hostIP)
			if err != nil {
				return poneglyph.Trace(err)
			}
		}
	case hostEntity.Remove:
		for _, hostIP := range req.Data {
			err = s.usecase.hostpool.RemoveHost(ctx, hostIP)
			if err != nil {
				return poneglyph.Trace(err)
			}
		}
	case hostEntity.HealthCheck:
		// execute health check
		var res hostEntity.HealthCheckAllHostResp
		res, err = m.HealthCheckAllHost(ctx)
		if err != nil {
			err = poneglyph.Trace(err)
		}

		// print health check result
		// Print the health check result
		fmt.Println("Forced Health Check Result:")
		fmt.Println("Healthy Hosts:")
		for _, host := range res.HealthyHosts {
			fmt.Printf("\t- %s\n", host)
		}

		fmt.Println("Down Hosts:")
		for _, host := range res.DownHosts {
			fmt.Printf("\t- %s\n", host)
		}
	default:
		return poneglyph.Trace(errors.New("unknown host operation"))
	}

	return nil
}

func (s Service) HealthCheckAllHost(ctx context.Context) (resp hostEntity.HealthCheckAllHostResp, err error) {
	// get all host
	var (
		hostMap     = make(map[string]bool)
		hostPoolMap = make(map[string]bool)
	)

	hostListFromPool, err := s.usecase.hostpool.GetHostListFromPool(ctx)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	// add all host to map
	for _, host := range hostListFromPool {
		hostMap[host] = true
		hostPoolMap[host] = true
	}

	for _, host := range s.config.HostList {
		hostMap[host] = true
	}

	// iterate host, try ping and check if it takes longer than a limit timeout(?)
	for host := range hostMap {
		isPingSuccess := s.usecase.poolclient.PingHost(ctx, host)
		if isPingSuccess {
			resp.HealthyHosts = append(resp.HealthyHosts, host)

			// if ok and not already in pool, add to pool
			isHostAlreadyInPool := hostPoolMap[host]
			if !isHostAlreadyInPool {
				errAddHost := s.usecase.hostpool.AddHost(ctx, host)
				if errAddHost != nil {
					// if fail, continue with other host but note the error
					errAddHost = poneglyph.Trace(err, fmt.Sprintf("add host fail for host: %s", host))
					fmt.Println(poneglyph.GetErrorLogMessage(errAddHost))
					continue
				}
			}
		} else {
			resp.DownHosts = append(resp.DownHosts, host)

			// if down/timeout, remove from host pool
			isHostAlreadyInPool := hostPoolMap[host]
			if !isHostAlreadyInPool {
				errRemoveHost := s.usecase.hostpool.RemoveHost(ctx, host)
				if errRemoveHost != nil {
					// if fail, continue with other host but note the error
					errRemoveHost = poneglyph.Trace(err, fmt.Sprintf("remove host fail for host: %s", host))
					fmt.Println(poneglyph.GetErrorLogMessage(errRemoveHost))
					continue
				}
			}
		}
	}

	return resp, nil
}
