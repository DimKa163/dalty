package protoerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(err error) error {
	internalError := mewInternalError(err)
	return status.Error(codes.Internal, internalError.Error())
}

func InvalidArgument(message string, args ...string) error {
	internalError := newInvalidArgsError(message, args)
	return status.Error(codes.Internal, internalError.Error())
}
