package service

import (
	"context"
	"go_crud/rpc_server/crud_pb"
)

func (*GrpcServer) Query(context.Context, *crud_pb.QueryRequest) (*crud_pb.QueryResponse, error) {

	return nil, nil
}
