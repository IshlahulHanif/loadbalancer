package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net"
	"net/http"
)

var (
	myIP string
)

func main() {
	var err error

	// init poneglyph settings
	poneglyph.SetProjectName("loadbalancer")
	poneglyph.SetIsPrintFromContentRoot(true)
	poneglyph.SetIsPrintFunctionName(true)
	poneglyph.SetIsPrintNewline(true)
	poneglyph.SetIsUseTabSeparator(false)

	myIP, err = GetLocalIP()
	if err != nil {
		err = poneglyph.Trace(err)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
	}

	http.HandleFunc("/", handleRequest)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		err = poneglyph.Trace(err)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		statusCode int
	)

	defer func() {
		if err != nil {
			HttpError(w, statusCode, err)
		}
	}()

	if r.Method != http.MethodPost {
		err = poneglyph.Trace(errors.New("method not allowed"))
		statusCode = http.StatusMethodNotAllowed
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload interface{}
	err = decoder.Decode(&payload)
	if err != nil {
		err = poneglyph.Trace(err, "Failed to decode body")
		statusCode = http.StatusBadRequest
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Host", myIP)

	// Respond with the received JSON payload
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		err = poneglyph.Trace(err, "Error encode body")
		statusCode = http.StatusInternalServerError
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HttpError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Println(poneglyph.GetErrorLogMessage(err))
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", poneglyph.Trace(err, "Failed to get local IP")
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}
