package forwarder

import (
	"bytes"
	"context"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/entity/forwarder"
	"io"
	"net/http"
)

func (s Service) ForwardRequest(ctx context.Context, req forwarder.ForwardRequestReq) (resp forwarder.ForwardRequestResp, err error) {
	resp.Header = map[string][]string{}

	// get hostpool
	host, err := s.usecase.hostpool.GetHost(ctx)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	// TODO: regulate http client request in a pkg, this function should only be getting the hostpool and forwarding the req, not building the resp etc
	// forward request by hitting the IP we got from the hostpool

	apiURL := host + req.Path

	// Add Query param
	if len(req.QueryParams) > 0 {
		apiURL += "?" + req.QueryParams.Encode()
	}

	// Create a new HTTP request to forward the request
	httpReq, err := http.NewRequest(req.Method, host+req.Path, bytes.NewBuffer(req.Body))
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	// Copy headers from the original request to the forward request
	for key, values := range req.Header {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	// TODO: after migrating http req to a pkg, add a feature to retry to another hostpool the request if the current request taking too long
	// Create an HTTP client and send the forward request
	client := http.DefaultClient
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}
	defer httpResp.Body.Close()

	// Read the response body from the forwarded request
	httpRespBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	for key, values := range httpResp.Header {
		for _, value := range values {
			resp.Header[key] = append(resp.Header[key], value)
		}
	}

	resp.Body = httpRespBody

	return resp, nil
}
