package service

import (
	"context"
	"go_crud/rpc_server/crud_pb"
)

func (*GrpcServer) Update(context.Context, *crud_pb.UpdateRequest) (*crud_pb.UpdateResponse, error) {

	return nil, nil
}
