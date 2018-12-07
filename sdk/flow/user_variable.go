package flow

import (
	"bufio"
	"fmt"
	"os"

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

// Set sets a variable from user input
func (s *UserVar) Set() {
	fmt.Printf("%s [%s]: ", s.Name, s.DefaultValue)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			s.set(line)
			return
		}
	}
	s.SetAsDefault()
}

// SetAsDefault sets the variable to its default value
func (s *UserVar) SetAsDefault() {
	s.set(s.DefaultValue)
}

// Set sets the value of the variable
func (s *UserVar) set(val string) {
	s.VarC <- val
	close(s.VarC)
}

// fanOutUserVars puts user variable setters on their own channel
func fanOutUserVars(incoming <-chan *grpcflow.Reply, flowsetterc chan<- *grpcflow.Request) (chan *grpcflow.Reply, <-chan *UserVar, chan error) {
	outgoingc := make(chan *grpcflow.Reply)
	userVarc := make(chan *UserVar)
	errc := make(chan error)

	go func() {
		defer func() {
			close(outgoingc)
			close(userVarc)
			close(errc)
			close(flowsetterc)
		}()

		for msg := range incoming {
			setRequest := &grpcflow.UserVariableSetter{}
			err := ptypes.UnmarshalAny(msg.Payload, setRequest)
			if err == nil {
				varC := make(chan string)
				userVar := &UserVar{
					VarC:         varC,
					Name:         setRequest.Name,
					DefaultValue: setRequest.Default,
				}

				userVarc <- userVar

				// Currently immediately sets the variable to its default value
				// TODO: pass setter along the chan
				inner, err := ptypes.MarshalAny(&grpcflow.UserVariableValue{
					Value: <-varC,
					Id:    setRequest.Id,
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

	return outgoingc, userVarc, errc
}
