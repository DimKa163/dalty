package server

import (
	"context"
	"github.com/DimKa163/dalty/internal/warehouse/core"
	"github.com/DimKa163/dalty/internal/warehouse/usecase"
	"github.com/DimKa163/dalty/pkg/api/warehouses/v1"

	"github.com/beevik/guid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PathServer struct {
	service *usecase.PathService
	warehousesv1.PathServiceServer
}

func NewPathServer(service *usecase.PathService) *PathServer {
	return &PathServer{
		service: service,
	}
}
func (ps *PathServer) Bind(server *grpc.Server) {
	warehousesv1.RegisterPathServiceServer(server, ps)
}
func (ps *PathServer) Get(ctx context.Context, in *warehousesv1.GetPath) (*warehousesv1.Path, error) {
	var protoPath warehousesv1.Path
	id, err := guid.ParseString(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	defWarehouse, err := guid.ParseString(in.GetDefaultWarehouseId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	path, err := ps.service.GetPath(ctx, id, defWarehouse)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	list := path.GetList()
	WarehouseCategorys := make([]*warehousesv1.Warehouse, len(list))
	for i, WarehouseCategory := range list {
		WarehouseCategorys[i] = mapWarehouseCategoryToProto(WarehouseCategory)
	}
	protoPath.SetNodes(WarehouseCategorys)
	return &protoPath, nil
}

func mapWarehouseCategoryToProto(WarehouseCategory *core.PathNode) *warehousesv1.Warehouse {
	var WarehouseCategoryProto warehousesv1.Warehouse
	WarehouseCategoryProto.SetId(WarehouseCategory.ID)
	it := WarehouseCategory.Value.(*core.Warehouse)
	WarehouseCategoryProto.SetName(it.Name)
	WarehouseCategoryProto.SetType(mapTypeToProto(it))
	if it.Info != nil {
		if it.Info.TimeZone != nil {
			WarehouseCategoryProto.SetTimeZone(it.Info.TimeZone.Code)
		}
		WarehouseCategoryProto.SetDescriptorGroup(it.Info.DescriptorGroup)
	}

	WarehouseCategoryProto.SetLevel(int32(WarehouseCategory.Level))
	WarehouseCategoryProto.SetAvailableRest(it.AvailableForBalance)
	WarehouseCategoryProto.SetAddress(it.Info.Address)
	WarehouseCategoryProto.SetOnlyStockPickupAllowed(it.OnlyStockPickupAllowed)

	return &WarehouseCategoryProto
}

func mapTypeToProto(n *core.Warehouse) warehousesv1.WarehouseType {
	switch n.Type {
	case core.WarehouseCategoryFree:
		return warehousesv1.WarehouseType_FREE
	case core.WarehouseCategoryMain:
		return warehousesv1.WarehouseType_MAIN
	case core.WarehouseCategoryCentral:
		return warehousesv1.WarehouseType_CENTRAL
	case core.WarehouseCategoryMall:
		return warehousesv1.WarehouseType_MALL
	case core.WarehouseCategoryTransit:
		return warehousesv1.WarehouseType_TRANSIT
	case core.WarehouseCategoryReservation:
		return warehousesv1.WarehouseType_RESERVATION
	case core.WarehouseCategoryLoses:
		return warehousesv1.WarehouseType_LOSES
	case core.WarehouseCategoryMarketing:
		return warehousesv1.WarehouseType_MARKETING
	case core.WarehouseCategoryExposition:
		return warehousesv1.WarehouseType_EXPOSITION
	case core.WarehouseCategoryPartner:
		return warehousesv1.WarehouseType_PARTNER
	case core.WarehouseCategoryPartner2:
		return warehousesv1.WarehouseType_PARTNER2
	case core.WarehouseCategoryFree2:
		return warehousesv1.WarehouseType_FREE2
	case core.WarehouseCategoryProblem:
		return warehousesv1.WarehouseType_PROBLEM
	case core.WarehouseCategoryRefund:
		return warehousesv1.WarehouseType_REFUND
	case core.WarehouseCategoryProduction:
		return warehousesv1.WarehouseType_PRODUCTION
	case core.WarehouseCategoryRecycling:
		return warehousesv1.WarehouseType_RECYCLING
	case core.WarehouseCategoryService:
		return warehousesv1.WarehouseType_SERVICE
	case core.WarehouseCategoryMaterial:
		return warehousesv1.WarehouseType_MATERIAL
	case core.WarehouseCategoryMarkdown:
		return warehousesv1.WarehouseType_MARKDOWN
	case core.WarehouseCategoryBuffer:
		return warehousesv1.WarehouseType_BUFFER
	case core.WarehouseCategoryDiscount:
		return warehousesv1.WarehouseType_DISCOUNT
	case core.WarehouseCategoryCentralMainIntermediate:
		return warehousesv1.WarehouseType_CENTRAL_MAIN_INTERMEDIATE
	case core.WarehouseCategoryMainCentraIntermediate:
		return warehousesv1.WarehouseType_MAIN_CENTRAL_INTERMEDIATE
	case core.WarehouseCategoryCentraFreeIntermediate:
		return warehousesv1.WarehouseType_CENTRAL_FREE_INTERMEDIATE
	case core.WarehouseCategoryFreeCentraIntermediate:
		return warehousesv1.WarehouseType_FREE_CENTRAL_INTERMEDIATE
	default:
		return warehousesv1.WarehouseType_UNSPECIFIED
	}
}
