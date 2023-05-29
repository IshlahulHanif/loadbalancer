package httpclient

import (
	"context"
	"crypto/tls"
	"net/url"
	"time"
)

type (
	Agent struct {
		Url     string
		Path    string
		Method  string
		Headers map[string][]string
		Param   url.Values
		Body    []byte
		//BodyJson     interface{}
		//IsJson       bool
		Timeout time.Duration
		//ResultStatus int
		Context context.Context
		TLS     *tls.Config
	}

	HttpRequestResponse struct {
		Body       []byte
		Header     map[string][]string
		StatusCode int
		IsSuccess  bool
	}
)
