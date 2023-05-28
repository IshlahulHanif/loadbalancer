package httpapi

type Module struct {
	service service
}

type service struct {
	forwarder forwarderMethod
	hostpool  hostpoolMethod
}
