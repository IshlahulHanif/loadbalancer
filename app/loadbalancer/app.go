package main

import (
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/api/httpapi"
	"github.com/loadbalancer/pkg/config"
	"net/http"
)

func main() {
	// init config
	conf := config.Config{
		HostList: []string{ //TODO: read from config file
			"http://127.0.0.1:8081",
			"http://127.0.0.1:8082",
			"http://127.0.0.1:8083",
		},
	}

	// init poneglyph settings
	poneglyph.SetProjectName("loadbalancer")
	poneglyph.SetIsPrintFromContentRoot(true)
	poneglyph.SetIsPrintFunctionName(true)
	poneglyph.SetIsPrintNewline(true)
	poneglyph.SetIsUseTabSeparator(false)

	// init http api
	httpApi, err := httpapi.GetInstance(conf)
	if err != nil {
		return
	}

	// register handlers
	http.HandleFunc("/", httpApi.HandlerForwardRequest)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		poneglyph.Trace(err)
		return
	}
}
