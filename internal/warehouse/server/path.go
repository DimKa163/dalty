package server

import (
	"context"

	"github.com/DimKa163/dalty/internal/warehouse/core"
	"github.com/DimKa163/dalty/internal/warehouse/proto"
	"github.com/DimKa163/dalty/internal/warehouse/usecase"
	"github.com/beevik/guid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PathServer struct {
	service *usecase.PathService
	proto.PathServiceServer
}

func NewPathServer(service *usecase.PathService) *PathServer {
	return &PathServer{
		service: service,
	}
}
func (ps *PathServer) Bind(server *grpc.Server) {
	proto.RegisterPathServiceServer(server, ps)
}
func (ps *PathServer) Get(ctx context.Context, in *proto.GetPath) (*proto.Path, error) {
	var protoPath proto.Path
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
	nodes := make([]*proto.Warehouse, len(list))
	for i, node := range list {
		nodes[i] = mapNodeToProto(node)
	}
	protoPath.SetNodes(nodes)
	return &protoPath, nil
}

func mapNodeToProto(node *core.PathNode) *proto.Warehouse {
	var nodeProto proto.Warehouse
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

func mapTypeToProto(n *core.Warehouse) proto.WarehouseType {
	switch n.Type {
	case core.NodeFree:
		return proto.WarehouseType_FREE
	case core.NodeMain:
		return proto.WarehouseType_MAIN
	case core.NodeCenter:
		return proto.WarehouseType_CENTRAL
	case core.NodeMall:
		return proto.WarehouseType_MALL
	case core.NodeTransit:
		return proto.WarehouseType_TRANSIT
	case core.NodeReservation:
		return proto.WarehouseType_RESERVATION
	case core.NodeLoses:
		return proto.WarehouseType_LOSES
	case core.NodeMarketing:
		return proto.WarehouseType_MARKETING
	case core.NodeExposition:
		return proto.WarehouseType_EXPOSITION
	case core.NodePartner:
		return proto.WarehouseType_PARTNER
	case core.NodePartner2:
		return proto.WarehouseType_PARTNER2
	case core.NodeFree2:
		return proto.WarehouseType_FREE2
	case core.NodeProblem:
		return proto.WarehouseType_PROBLEM
	case core.NodeRefund:
		return proto.WarehouseType_REFUND
	case core.NodeProduction:
		return proto.WarehouseType_PRODUCTION
	case core.NodeRecycling:
		return proto.WarehouseType_RECYCLING
	case core.NodeService:
		return proto.WarehouseType_SERVICE
	case core.NodeMaterial:
		return proto.WarehouseType_MATERIAL
	case core.NodeMarkdown:
		return proto.WarehouseType_MARKDOWN
	case core.NodeBuffer:
		return proto.WarehouseType_BUFFER
	case core.NodeDiscount:
		return proto.WarehouseType_DISCOUNT
	case core.NodeCentralMainIntermediate:
		return proto.WarehouseType_CENTRAL_MAIN_INTERMEDIATE
	case core.NodeMainCentralIntermediate:
		return proto.WarehouseType_MAIN_CENTRAL_INTERMEDIATE
	case core.NodeCentralFreeIntermediate:
		return proto.WarehouseType_CENTRAL_FREE_INTERMEDIATE
	case core.NodeFreeCentralIntermediate:
		return proto.WarehouseType_FREE_CENTRAL_INTERMEDIATE
	default:
		return proto.WarehouseType_UNRECOGNIZED
	}
}
