package host

type (
	ManageHostReq struct {
		Operation Operation `json:"operation"`
		Data      []string  `json:"data"`
	}

	HealthCheckAllHostResp struct {
		HealthyHosts []string
		DownHosts    []string
	}
)
