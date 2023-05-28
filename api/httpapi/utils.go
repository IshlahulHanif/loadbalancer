package httpapi

import (
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net/http"
)

func HttpError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Println(poneglyph.GetErrorLogMessage(err))
}
