package server

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/internal/rest/usecase"
	restsv1 "github.com/DimKa163/dalty/pkg/api/rests/v1"
	"github.com/DimKa163/dalty/pkg/daltyerrors"
	"github.com/DimKa163/dalty/pkg/daltyerrors/protoerr"
	"github.com/beevik/guid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RestServer struct {
	service *usecase.RestService
	restsv1.UnimplementedRestServiceServer
}

func NewRestServer(service *usecase.RestService) *RestServer {
	return &RestServer{
		service: service,
	}
}

func (s *RestServer) Bind(server *grpc.Server) {
	restsv1.RegisterRestServiceServer(server, s)
}

func (s *RestServer) Get(ctx context.Context, in *restsv1.RestRequest) (*restsv1.RestResponse, error) {
	var result restsv1.RestResponse
	req, err := toIn(in)
	if err != nil {
		return nil, err
	}
	res, err := s.service.GetRest(ctx, req)
	if err != nil {
		var daltyErr *daltyerrors.DaltyError
		if errors.As(err, &daltyErr) {
			return nil, protoerr.Handle(daltyErr)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result.SetBalanceMap(toOut(res))
	return &result, nil
}

func toIn(in *restsv1.RestRequest) (*usecase.RestRequest, error) {
	var req usecase.RestRequest
	if !in.HasFilialId() {
		return nil, protoerr.InvalidArgument("filia_id is empty", &protoerr.ValidationError{
			Message: "undefined filial_id",
			Members: []string{
				"filial_id",
			},
		})
	}
	filialID, err := guid.ParseString(in.GetFilialId())
	if err != nil {
		return nil, protoerr.InvalidArgument(err.Error(), &protoerr.ValidationError{
			Message: "invalid filial_id",
			Members: []string{
				"filial_id",
			},
		})
	}
	req.FilialID = *filialID
	subWarehousesIDS := make([]guid.Guid, len(in.GetWarehouseIds()))
	for i, wh := range in.GetWarehouseIds() {
		id, err := guid.ParseString(wh)
		if err != nil {
			return nil, protoerr.InvalidArgument(err.Error(), &protoerr.ValidationError{
				Message: "invalid warehouse_id",
				Members: []string{
					"warehouse_ids",
				},
			})
		}
		subWarehousesIDS[i] = *id
	}
	productIDS := make([]guid.Guid, len(in.GetProductIds()))
	for i, pid := range in.GetProductIds() {
		id, err := guid.ParseString(pid)
		if err != nil {
			return nil, protoerr.InvalidArgument(err.Error(), &protoerr.ValidationError{
				Message: "invalid product_id",
				Members: []string{
					"product_ids",
				},
			})
		}
		productIDS[i] = *id
	}
	req.WarehousesIDS = subWarehousesIDS
	req.ProductIDS = productIDS
	return &req, nil
}

func toOut(out *usecase.RestResult) map[string]*restsv1.RestMap {
	balanceMap := make(map[string]*restsv1.RestMap)
	for i, v := range out.ProductMap {
		var m restsv1.RestMap
		mp := make(map[string]*restsv1.Rest)
		for j, v1 := range v {
			var r restsv1.Rest
			r.SetProductId(v1.ProductID.String())
			r.SetWarehouseId(v1.SubWarehouseID.String())
			r.SetFilialId(v1.FilialID.String())
			r.SetQuantity(v1.Quantity.InexactFloat64())
			mp[j.String()] = &r
		}
		m.SetRests(mp)
		balanceMap[i.String()] = &m
	}
	return balanceMap
}
