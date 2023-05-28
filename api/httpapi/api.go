package httpapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/entity/forwarder"
	"io"
	"net/http"
)

func (m Module) HandlerForwardRequest(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = context.Background()
		err error
	)

	if r.Method != http.MethodPost {
		// TODO: these error handling can be extracted as a function
		err = poneglyph.Trace(errors.New("unsupported HTTP method"))
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(poneglyph.GetErrorLogMessage(err))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = poneglyph.Trace(err, "Error reading request body")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(poneglyph.GetErrorLogMessage(err))
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

	r.URL.String()

	resp, err := m.service.forwarder.ForwardRequest(ctx, req)
	if err != nil {
		err = poneglyph.Trace(err)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(poneglyph.GetErrorLogMessage(err))
		return
	}

	// Set headers
	for key, values := range resp.Header {
		for _, value := range values {
			resp.Header[key] = append(resp.Header[key], value)
			w.Header().Set(key, value)
		}
	}

	//err = json.NewEncoder(w).Encode(resp.Body) // TODO: should we use json encoder instead?
	_, err = fmt.Fprint(w, string(resp.Body))
	if err != nil {
		err = poneglyph.Trace(err)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(poneglyph.GetErrorLogMessage(err))
		return
	}
}
