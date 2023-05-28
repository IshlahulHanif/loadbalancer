package poolclient

import (
	"context"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
	"net"
	"strings"
)

func (u Usecase) PingHost(ctx context.Context, host string) (isPingSuccess bool) {
	// TODO: improve ping logics, maybe add reason etc
	timeout := u.config.PingTimeout

	// TODO: better handling for weird cases because there is http:// or not
	// handle containing http case
	hostWithoutHttpSlice := strings.Split(host, "//")
	var hostWithoutHttp string
	if len(hostWithoutHttpSlice) > 1 {
		hostWithoutHttp = hostWithoutHttpSlice[1]
	}

	_, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		err = poneglyph.Trace(err)
		fmt.Println(poneglyph.GetErrorLogMessage(err))

		// try without http://
		_, err = net.DialTimeout("tcp", hostWithoutHttp, timeout)
		if err != nil {
			err = poneglyph.Trace(err)
			fmt.Println(poneglyph.GetErrorLogMessage(err))

			return false
		}
	}

	return true
}
