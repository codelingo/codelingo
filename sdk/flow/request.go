package flow

import (
	"context"

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

	rpcReq := &grpcflow.Request{
		Flow:    flowName,
		Payload: payload,
	}

	replyc, runErrc, err := Request(ctx, rpcReq)
	if err != nil {
		return nil, nil, nil, errors.Trace(err)
	}

	itemc, marshalErrc := MarshalChan(replyc, newItem)
	return itemc, ErrFanIn(runErrc, marshalErrc), cancel, nil
}

func Request(ctx context.Context, req *grpcflow.Request) (chan *grpcflow.Reply, chan error, error) {
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

	replyc, runErrc, err := c.Run(ctx, req)
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
