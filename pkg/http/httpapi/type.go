package httpapi

type Module struct {
	usecase usecase
}

type usecase struct {
	loadbalancer loadbalancerMethod
}
