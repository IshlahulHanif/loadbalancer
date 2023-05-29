package forwarder

import "net/url"

type (
	ForwardRequestReq struct {
		Body        []byte
		Header      map[string][]string
		Path        string
		QueryParams url.Values
		Method      string
	}

	ForwardRequestResp struct {
		Body       []byte
		Header     map[string][]string
		StatusCode int
	}
)
