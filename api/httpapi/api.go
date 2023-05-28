package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/entity/forwarder"
	"github.com/loadbalancer/entity/host"
	"io"
	"net/http"
)

func (m Module) HandlerForwardRequest(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = context.Background()
		err        error
		statusCode int
	)

	defer func() {
		if err != nil {
			HttpError(w, statusCode, err)
		}
	}()

	if r.Method != http.MethodPost {
		err = poneglyph.Trace(errors.New("unsupported HTTP method"))
		statusCode = http.StatusMethodNotAllowed
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = poneglyph.Trace(err, "Error reading request body")
		statusCode = http.StatusBadRequest
		return
	}
	defer r.Body.Close()

	req := forwarder.ForwardRequestReq{
		Body:        body,
		Header:      r.Header,
		Path:        r.URL.Path,
		QueryParams: r.URL.Query(),
		Method:      r.Method,
	}

	resp, err := m.service.forwarder.ForwardRequest(ctx, req)
	if err != nil {
		err = poneglyph.Trace(err)
		statusCode = http.StatusInternalServerError
		return
	}

	// Set headers
	for key, values := range resp.Header {
		for _, value := range values {
			resp.Header[key] = append(resp.Header[key], value)
			w.Header().Set(key, value)
		}
	}

	// We use raw encoder because this api purpose is to forward anything
	_, err = fmt.Fprint(w, string(resp.Body))
	if err != nil {
		err = poneglyph.Trace(err)
		statusCode = http.StatusInternalServerError
		return
	}
}

func (m Module) HandlerManageHostRequest(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = context.Background()
		err        error
		statusCode int
	)

	defer func() {
		if err != nil {
			HttpError(w, statusCode, err)
		}
	}()

	if r.Method != http.MethodPost {
		err = poneglyph.Trace(errors.New("unsupported HTTP method"))
		statusCode = http.StatusMethodNotAllowed
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload host.ManageHostReq
	err = decoder.Decode(&payload)
	if err != nil {
		err = poneglyph.Trace(err, "Failed to decode body")
		statusCode = http.StatusBadRequest
		return
	}

	err = m.service.hostpool.ManageHost(ctx, payload)
	if err != nil {
		err = poneglyph.Trace(err)
		statusCode = http.StatusInternalServerError
		return
	}

	result := map[string]string{
		"result": "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		err = poneglyph.Trace(err, "Error encode body")
		statusCode = http.StatusInternalServerError
		return
	}

	w.WriteHeader(http.StatusOK)
}
