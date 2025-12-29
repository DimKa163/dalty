package server

import (
	"context"

	"github.com/DimKa163/dalty/internal/graph/core"
	"github.com/DimKa163/dalty/internal/graph/proto"
	"github.com/DimKa163/dalty/internal/graph/usecase"
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
	nodes := make([]*proto.Node, len(list))
	for i, node := range list {
		nodes[i] = mapNodeToProto(node)
	}
	protoPath.SetNodes(nodes)
	return &protoPath, nil
}

func mapNodeToProto(node *core.PathNode) *proto.Node {
	var nodeProto proto.Node
	nodeProto.SetId(node.ID.String())
	nodeProto.SetName(node.Name)
	nodeProto.SetType(mapTypeToProto(node))
	nodeProto.SetTimeZone(node.Code)
	nodeProto.SetLevel(int32(node.Level))
	nodeProto.SetAvailableRest(node.AvailableRest)
	nodeProto.SetAddress(node.Address)
	nodeProto.SetOnlyStockPickupAllowed(node.OnlyStockPickupAllowed)
	nodeProto.SetDescriptorGroup(node.DescriptorGroup)
	return &nodeProto
}

func mapTypeToProto(n *core.PathNode) proto.NodeType {
	switch n.Type {
	case core.NodeFree:
		return proto.NodeType_FREE
	case core.NodeMain:
		return proto.NodeType_MAIN
	case core.NodeCenter:
		return proto.NodeType_CENTRAL
	case core.NodeMall:
		return proto.NodeType_MALL
	case core.NodeTransit:
		return proto.NodeType_TRANSIT
	case core.NodeReservation:
		return proto.NodeType_RESERVATION
	case core.NodeLoses:
		return proto.NodeType_LOSES
	case core.NodeMarketing:
		return proto.NodeType_MARKETING
	case core.NodeExposition:
		return proto.NodeType_EXPOSITION
	case core.NodePartner:
		return proto.NodeType_PARTNER
	case core.NodePartner2:
		return proto.NodeType_PARTNER2
	case core.NodeFree2:
		return proto.NodeType_FREE2
	case core.NodeProblem:
		return proto.NodeType_PROBLEM
	case core.NodeRefund:
		return proto.NodeType_REFUND
	case core.NodeProduction:
		return proto.NodeType_PRODUCTION
	case core.NodeRecycling:
		return proto.NodeType_RECYCLING
	case core.NodeService:
		return proto.NodeType_SERVICE
	case core.NodeMaterial:
		return proto.NodeType_MATERIAL
	case core.NodeMarkdown:
		return proto.NodeType_MARKDOWN
	case core.NodeBuffer:
		return proto.NodeType_BUFFER
	case core.NodeDiscount:
		return proto.NodeType_DISCOUNT
	case core.NodeCentralMainIntermediate:
		return proto.NodeType_CENTRAL_MAIN_INTERMEDIATE
	case core.NodeMainCentralIntermediate:
		return proto.NodeType_MAIN_CENTRAL_INTERMEDIATE
	case core.NodeCentralFreeIntermediate:
		return proto.NodeType_CENTRAL_FREE_INTERMEDIATE
	case core.NodeFreeCentralIntermediate:
		return proto.NodeType_FREE_CENTRAL_INTERMEDIATE
	default:
		return proto.NodeType_UNRECOGNIZED
	}
}
