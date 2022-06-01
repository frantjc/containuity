package svc

import "google.golang.org/grpc"

type Service interface {
	Register(grpc.ServiceRegistrar)
}
