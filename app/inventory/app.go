package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net"
	"net/http"
	"time"
)

var (
	myIP        string
	delaySecond time.Duration
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

	// Create two routers for different ports
	requestRouter := http.NewServeMux()
	settingsRouter := http.NewServeMux()

	// register handlers
	requestRouter.HandleFunc("/", handleBounceRequest)
	settingsRouter.HandleFunc("/bounce/delay/add", handleAddDelay)

	var errChan = make(chan error)

	// Start the requestRouter server on port 8080
	go func() {
		errServer := http.ListenAndServe(":8080", requestRouter)
		if errServer != nil {
			errChan <- poneglyph.Trace(errServer)
		}
	}()

	// Start the hostManagementRouter server on port 9090
	go func() {
		errServer := http.ListenAndServe(":9090", settingsRouter)
		if errServer != nil {
			errChan <- poneglyph.Trace(errServer)
		}
	}()

	err = <-errChan
	if err != nil {
		err = poneglyph.Trace(err)
		fmt.Println(poneglyph.GetLogErrorTrace(err))
		// TODO: consume all errChan to avoid memory leak
	}

	// TODO: make sure to Shutdown the routers gracefully
}

func handleAddDelay(w http.ResponseWriter, r *http.Request) {
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
	var payload struct {
		DelaySecond int `json:"delay_second"`
	}

	err = decoder.Decode(&payload)
	if err != nil {
		err = poneglyph.Trace(err, "Failed to decode body")
		statusCode = http.StatusBadRequest
		return
	}

	delaySecond = time.Duration(payload.DelaySecond) * time.Second

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

func handleBounceRequest(w http.ResponseWriter, r *http.Request) {
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

	if delaySecond > 0 {
		time.Sleep(delaySecond)
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
