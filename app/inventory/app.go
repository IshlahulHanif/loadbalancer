package main

import (
	"encoding/json"
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload map[string]interface{}
	err := decoder.Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: use the pattern from the round robin app
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Host", myIP)

	// Respond with the received JSON payload
	responseJSON, _ := json.Marshal(payload)
	_, err = fmt.Fprint(w, string(responseJSON))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
