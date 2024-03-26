package service

import (
	"context"
	"go_crud/rpc_server/crud_pb"
)

func (*GrpcServer) Delete(c context.Context, req *crud_pb.DeleteRequest) (res *crud_pb.DeleteResponse, e error) {

	return nil, nil
}
