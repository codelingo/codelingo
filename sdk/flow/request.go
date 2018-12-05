package flow

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"

	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/service"
	grpcclient "github.com/codelingo/lingo/service/grpc"
	grpcflow "github.com/codelingo/rpc/flow"
	"github.com/codelingo/rpc/flow/client"
	"github.com/juju/errors"
)

func RunFlow(flowName string, req proto.Message, newItem func() proto.Message) (chan proto.Message, chan error, func(), error) {

	ctx, cancel := util.UserCancelContext(context.Background())

	payload, err := ptypes.MarshalAny(req)
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	rpcReqC := make(chan *grpcflow.Request)
	go func() {
		rpcReqC <- &grpcflow.Request{
			Flow:    flowName,
			Payload: payload,
		}
	}()

	allReplyc, runErrc, err := Request(ctx, rpcReqC)
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	// TODO: send setter chan higher
	replyc, _, setterErrc := SplitSetters(allReplyc, rpcReqC)

	itemc, marshalErrc := MarshalChan(replyc, newItem)
	return itemc, ErrFanIn(ErrFanIn(runErrc, marshalErrc), setterErrc), cancel, nil
}

func Request(ctx context.Context, reqC <-chan *grpcflow.Request) (chan *grpcflow.Reply, chan error, error) {
	conn, err := service.GrpcConnection(service.LocalClient, service.FlowServer, true)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	c := client.NewFlowClient(conn)

	// Create context with metadata
	ctx, err = grpcclient.AddUsernameToCtx(ctx)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	replyc, runErrc, err := c.Run(ctx, reqC, nil)
	return replyc, runErrc, errors.Trace(err)
}

func MarshalChan(replyc chan *grpcflow.Reply, newItem func() proto.Message) (chan proto.Message, chan error) {
	itemc := make(chan proto.Message)
	errc := make(chan error)

	go func() {
		for reply := range replyc {
			if reply.IsHeartbeat {
				continue
			}
			if reply.Error != "" {
				errc <- errors.New(reply.Error)
				continue
			}

			util.Logger.Debug("got reply %v", reply)

			item := newItem()
			err := ptypes.UnmarshalAny(reply.Payload, item)
			if err != nil {
				errc <- err
				continue
			}

			itemc <- item
		}
		close(errc)
		close(itemc)
	}()

	return itemc, errc
}

// SplitSetters puts user variable setters on their own channel
func SplitSetters(incoming <-chan *grpcflow.Reply, flowsetterc chan *grpcflow.Request) (chan *grpcflow.Reply, <-chan *Setter, chan error) {
	outgoingc := make(chan *grpcflow.Reply)
	clientsetterc := make(chan *Setter)
	errc := make(chan error)

	go func() {
		defer close(outgoingc)
		defer close(clientsetterc)
		defer close(errc)

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

				fmt.Println("setting value of", setter.Name, "to", setter.Default)
				flowsetterc <- &grpcflow.Request{
					Payload: inner,
				}
			} else {
				fmt.Println("GOT ERROR")
				outgoingc <- msg
			}
		}
	}()

	return outgoingc, clientsetterc, errc
}

// A Setter allows users and other external agents to set variable values while a query
// is being executed.
// TODO: setter code is copied from the Platform
type Setter struct {
	VarC         chan<- string
	Name         string
	DefaultValue string
}

// SetAsDefault sets the variable to its default value
func (s *Setter) SetAsDefault() {
	s.VarC <- s.DefaultValue
	close(s.VarC)
}
