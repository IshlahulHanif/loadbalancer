package httpclient

import (
	"bytes"
	"github.com/IshlahulHanif/poneglyph"
	"io"
	"net/http"
	"net/url"
)

func NewHTTPAgent() *Agent {
	return &Agent{
		Headers: make(map[string][]string),
	}
}

func (a *Agent) DoReq() (resp HttpRequestResponse, err error) {
	// TODO: create a better default response data
	resp.Header = make(map[string][]string)

	// build url
	apiUrl, err := url.Parse(a.Url)
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	apiUrl.Path = a.Path
	apiUrl.RawQuery = a.Param.Encode()

	httpReq, err := http.NewRequest(a.Method, apiUrl.String(), bytes.NewBuffer(a.Body))
	if err != nil {
		return resp, poneglyph.Trace(err)
	}

	// Copy headers from the original request to the forward request
	for key, values := range a.Headers {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	defaultHttpClient := http.DefaultClient
	defaultHttpClient.Timeout = a.Timeout

	httpResp, err := defaultHttpClient.Do(httpReq)
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
	resp.IsSuccess = true
	resp.StatusCode = httpResp.StatusCode

	return resp, nil
}
