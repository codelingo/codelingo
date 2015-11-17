package service

import "google.golang.org/grpc"

type Service interface {
	Start() error
	Stop() error
	IsRunning() bool
	DialGRPC() (*grpc.ClientConn, error)
}
