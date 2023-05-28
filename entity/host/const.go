package host

type Operation int

const (
	Unknown     Operation = -1
	Add         Operation = 1
	Remove      Operation = 2
	HealthCheck Operation = 3
)
