package hostpool

type Usecase struct {
	roundRobinQueue RoundRobinQueue
}

type RoundRobinQueue struct {
	queue []string
	index int
}
