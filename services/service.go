package services

import "google.golang.org/grpc"

type Service interface {
	Client() (interface{}, error)
	Register(grpc.ServiceRegistrar)
}
