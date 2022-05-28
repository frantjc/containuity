package services

import (
	"github.com/frantjc/sequence/runtime"
	"google.golang.org/grpc"
)

type service struct {
	runtime runtime.Runtime
}

type Service interface {
	Register(grpc.ServiceRegistrar)
}
