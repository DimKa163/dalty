package protoerr

import (
	"github.com/DimKa163/dalty/pkg/api/common/v1"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

type ValidationError struct {
	Message string
	Members []string
}

func InvalidArgument(message string, errs ...*ValidationError) error {
	st := status.New(codes.InvalidArgument, message)
	itemErrors := make([]*commonv1.ValidationError, len(errs))
	var detail commonv1.ErrorDetail
	detail.SetCode(001)
	detail.SetMessage(message)
	for i, err := range errs {
		var itemError commonv1.ValidationError
		itemError.SetMessage(err.Message)
		itemError.SetMembers(err.Members)
		itemErrors[i] = &itemError
	}
	detail.SetValidationErrors(itemErrors)
	st, _ = st.WithDetails(&detail)
	return st.Err()
}

type EntityError struct {
	Message                string
	Nrec                   string
	IntegrationID          string
	DestinationWarehouseID string
	ShipmentWarehouseID    string
}

func NotFound(message string, errs []*EntityError) error {
	st := status.New(codes.NotFound, message)
	var detail commonv1.ErrorDetail
	detail.SetCode(001)
	detail.SetMessage(message)
	return st.Err()
}

func Handle(err *daltyerrors.DaltyError) error {
	var st *status.Status
	switch err.Type {
	case daltyerrors.DaltyErrorTypeInvalidRequest:
		st = status.New(codes.InvalidArgument, err.Error())
	case daltyerrors.DaltyErrorTypeResourceNotFound:
		st = status.New(codes.NotFound, err.Error())
	case daltyerrors.DaltyErrorTypeBusinessError:
		st = status.New(codes.FailedPrecondition, err.Error())
	default:
		st = status.New(codes.Internal, err.Error())
	}
	entErros := make([]*commonv1.EntityError, len(err.EntityErrors))
	for i, entErr := range err.EntityErrors {
		var itemError commonv1.EntityError
		itemError.SetId(entErr.ID)
		itemError.SetEntityName(entErr.EntityName)
		entErros[i] = &itemError
	}
	var detail commonv1.ErrorDetail
	detail.SetCode(int32(err.Code))
	detail.SetMessage(err.Message)
	detail.SetEntityErrors(entErros)
	st, _ = st.WithDetails(&detail)
	return st.Err()
}

func Flatten(status *status.Status) []*commonv1.ErrorDetail {
	details := make([]*commonv1.ErrorDetail, 0, len(status.Details()))
	for _, detail := range status.Details() {
		errDetail, ok := detail.(*commonv1.ErrorDetail)
		if !ok {
			continue
		}
		details = append(details, errDetail)
	}
	return details

}
