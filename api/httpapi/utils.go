package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net/http"
)

func HttpError(w http.ResponseWriter, statusCode int, err error) {
	// TODO: improve error message template
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	// write error message
	var errMsg = map[string]string{
		"error": err.Error(),
	}

	errByte, _ := json.Marshal(errMsg)

	fmt.Fprint(w, string(errByte))

	fmt.Println(poneglyph.GetErrorLogMessage(err))
}
