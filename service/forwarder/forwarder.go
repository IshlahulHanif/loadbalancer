package forwarder

import (
	"context"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/entity/forwarder"
	"github.com/loadbalancer/pkg/http/httpclient"
	"strings"
)

func (s Service) ForwardRequest(ctx context.Context, req forwarder.ForwardRequestReq) (resp forwarder.ForwardRequestResp, err error) {
	resp.Header = map[string][]string{}

	// get hostpool
	host, err := s.usecase.hostpool.GetHostDequeueRoundRobin(ctx)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	// create a request
	agent := httpclient.NewHTTPAgent()
	agent.Url = host
	agent.Path = req.Path
	agent.Param = req.QueryParams
	agent.Body = req.Body
	agent.Headers = req.Header
	agent.Method = req.Method
	agent.Timeout = s.config.RequestTimeout

	httpResp, err := agent.DoReq()
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			// if timeout, remove host from pool //TODO: should only remove host from pool after multiple timeout, can use circuit breaker
			errRemoveHost := s.usecase.hostpool.RemoveHost(ctx, host)
			if errRemoveHost != nil {
				errRemoveHost = poneglyph.Trace(errRemoveHost)
				fmt.Println(poneglyph.GetErrorLogMessage(errRemoveHost))
			}
		}

		return resp, poneglyph.Trace(err)
	}

	// loop try if its time out and try the next host

	resp.Body = httpResp.Body
	resp.StatusCode = httpResp.StatusCode
	resp.Header = httpResp.Header

	return resp, nil
}
