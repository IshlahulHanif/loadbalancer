package hostpool

type (
	Service struct {
		usecase usecase
	}

	usecase struct {
		hostpool UsecaseHostpool
	}
)
