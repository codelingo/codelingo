package flow

import (
	grpcflow "github.com/codelingo/rpc/flow"
	"github.com/golang/protobuf/ptypes"
	"github.com/juju/errors"
)

// A UserVar allows users and other external agents to set variable values while a query
// is being executed.
// TODO: UserVar code is copied from the Platform
type UserVar struct {
	VarC         chan<- string
	Name         string
	DefaultValue string
}

// Set sets the value of the variable
func (s *UserVar) Set(val string) {
	s.VarC <- val
	close(s.VarC)
}

// SetAsDefault sets the variable to its default value
func (s *UserVar) SetAsDefault() {
	s.Set(s.DefaultValue)
}

// fanOutUserVars puts user variable setters on their own channel
func fanOutUserVars(incoming <-chan *grpcflow.Reply, flowsetterc chan<- *grpcflow.Request) (chan *grpcflow.Reply, <-chan *UserVar, chan error) {
	outgoingc := make(chan *grpcflow.Reply)
	clientsetterc := make(chan *UserVar)
	errc := make(chan error)

	go func() {
		defer func() {
			close(outgoingc)
			close(clientsetterc)
			close(errc)
			close(flowsetterc)
		}()

		for msg := range incoming {
			setter := &grpcflow.UserVariableSetter{}
			err := ptypes.UnmarshalAny(msg.Payload, setter)
			if err == nil {
				// Currently immediately sets the variable to its default value
				// TODO: pass setter along the chan
				inner, err := ptypes.MarshalAny(&grpcflow.UserVariableValue{
					Value: setter.Default,
					Id:    setter.Id,
				})
				if err != nil {
					errc <- errors.Trace(err)
				}

				flowsetterc <- &grpcflow.Request{
					Payload: inner,
				}
			} else {
				outgoingc <- msg
			}
		}
	}()

	return outgoingc, clientsetterc, errc
}
