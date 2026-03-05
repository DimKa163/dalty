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
	nodes := make([]*warehousesv1.Warehouse, len(list))
	for i, node := range list {
		nodes[i] = mapNodeToProto(node)
	}
	protoPath.SetNodes(nodes)
	return &protoPath, nil
}

func mapNodeToProto(node *core.PathNode) *warehousesv1.Warehouse {
	var nodeProto warehousesv1.Warehouse
	nodeProto.SetId(node.ID)
	it := node.Value.(*core.Warehouse)
	nodeProto.SetName(it.Name)
	nodeProto.SetType(mapTypeToProto(it))
	if it.Info != nil {
		if it.Info.TimeZone != nil {
			nodeProto.SetTimeZone(it.Info.TimeZone.Code)
		}
		nodeProto.SetDescriptorGroup(it.Info.DescriptorGroup)
	}

	nodeProto.SetLevel(int32(node.Level))
	nodeProto.SetAvailableRest(it.AvailableForBalance)
	nodeProto.SetAddress(it.Info.Address)
	nodeProto.SetOnlyStockPickupAllowed(it.OnlyStockPickupAllowed)

	return &nodeProto
}

func mapTypeToProto(n *core.Warehouse) warehousesv1.WarehouseType {
	switch n.Type {
	case core.NodeFree:
		return warehousesv1.WarehouseType_FREE
	case core.NodeMain:
		return warehousesv1.WarehouseType_MAIN
	case core.NodeCenter:
		return warehousesv1.WarehouseType_CENTRAL
	case core.NodeMall:
		return warehousesv1.WarehouseType_MALL
	case core.NodeTransit:
		return warehousesv1.WarehouseType_TRANSIT
	case core.NodeReservation:
		return warehousesv1.WarehouseType_RESERVATION
	case core.NodeLoses:
		return warehousesv1.WarehouseType_LOSES
	case core.NodeMarketing:
		return warehousesv1.WarehouseType_MARKETING
	case core.NodeExposition:
		return warehousesv1.WarehouseType_EXPOSITION
	case core.NodePartner:
		return warehousesv1.WarehouseType_PARTNER
	case core.NodePartner2:
		return warehousesv1.WarehouseType_PARTNER2
	case core.NodeFree2:
		return warehousesv1.WarehouseType_FREE2
	case core.NodeProblem:
		return warehousesv1.WarehouseType_PROBLEM
	case core.NodeRefund:
		return warehousesv1.WarehouseType_REFUND
	case core.NodeProduction:
		return warehousesv1.WarehouseType_PRODUCTION
	case core.NodeRecycling:
		return warehousesv1.WarehouseType_RECYCLING
	case core.NodeService:
		return warehousesv1.WarehouseType_SERVICE
	case core.NodeMaterial:
		return warehousesv1.WarehouseType_MATERIAL
	case core.NodeMarkdown:
		return warehousesv1.WarehouseType_MARKDOWN
	case core.NodeBuffer:
		return warehousesv1.WarehouseType_BUFFER
	case core.NodeDiscount:
		return warehousesv1.WarehouseType_DISCOUNT
	case core.NodeCentralMainIntermediate:
		return warehousesv1.WarehouseType_CENTRAL_MAIN_INTERMEDIATE
	case core.NodeMainCentralIntermediate:
		return warehousesv1.WarehouseType_MAIN_CENTRAL_INTERMEDIATE
	case core.NodeCentralFreeIntermediate:
		return warehousesv1.WarehouseType_CENTRAL_FREE_INTERMEDIATE
	case core.NodeFreeCentralIntermediate:
		return warehousesv1.WarehouseType_FREE_CENTRAL_INTERMEDIATE
	default:
		return warehousesv1.WarehouseType_UNRECOGNIZED
	}
}
