package service

import (
	"context"
	"go_crud/rpc_server/crud_pb"
)

func (*GrpcServer) Delete(context.Context, *crud_pb.DeleteRequest) (*crud_pb.DeleteResponse, error) {

	return nil, nil
}
