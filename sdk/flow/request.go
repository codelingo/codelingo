package flow

import (
	"context"

	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/service"
	grpcclient "github.com/codelingo/lingo/service/grpc"
	grpcflow "github.com/codelingo/rpc/flow"
	"github.com/codelingo/rpc/flow/client"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/juju/errors"
)

// RunFlow calls a given flow on the flow server and marshalls replies as a given type
func RunFlow(flowName string, req proto.Message, newItem func() proto.Message, setDefaults func(proto.Message) proto.Message) (chan proto.Message, <-chan *UserVar, chan error, func(), error) {
	ctx, cancel := util.UserCancelContext(context.Background())

	payload, err := ptypes.MarshalAny(req)
	if err != nil {
		return nil, nil, nil, nil, errors.Trace(err)
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
		return nil, nil, nil, nil, errors.Trace(err)
	}

	replyc, userVarC, setterErrc := fanOutUserVars(allReplyc, rpcReqC)
	itemc, marshalErrc := MarshalChan(replyc, newItem, setDefaults)
	return itemc, userVarC, ErrFanIn(ErrFanIn(runErrc, marshalErrc), setterErrc), cancel, nil
}

func Request(ctx context.Context, reqC <-chan *grpcflow.Request) (chan *grpcflow.Reply, chan error, error) {
	util.Logger.Debug("opening connection to flow server ...")
	conn, err := service.GrpcConnection(service.LocalClient, service.FlowServer, true)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	util.Logger.Debug("...connection to flow server opened")

	c := client.NewFlowClient(conn)

	// Create context with metadata
	ctx, err = grpcclient.AddUsernameToCtx(ctx)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	replyc, runErrc, err := c.Run(ctx, reqC)
	return replyc, runErrc, errors.Trace(err)
}

func MarshalChan(replyc chan *grpcflow.Reply, newItem func() proto.Message, setDefaults func(proto.Message) proto.Message) (chan proto.Message, chan error) {
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

			itemc <- setDefaults(item)
		}
		close(errc)
		close(itemc)
	}()

	return itemc, errc
}
