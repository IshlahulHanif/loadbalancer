package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net/http"
)

func (m Module) HandlerBouncerApi(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = context.Background()
		err error
	)

	if r.Method != http.MethodPost {
		// TODO: these error handling can be extracted as a function
		err = poneglyph.Trace(errors.New("unsupported HTTP method"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload map[string]interface{} // make sure at least the payload is somewhat a valid json
	err = decoder.Decode(&payload)
	if err != nil {
		err = poneglyph.Trace(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
		return
	}

	resp, err := m.usecase.loadbalancer.BouncerApi(ctx, payload)
	if err != nil {
		err = poneglyph.Trace(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		err = poneglyph.Trace(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
