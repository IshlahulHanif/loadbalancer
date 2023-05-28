package forwarder

type (
	Service struct {
		usecase usecase
	}

	usecase struct {
		hostpool UsecaseHostpool
	}
)
