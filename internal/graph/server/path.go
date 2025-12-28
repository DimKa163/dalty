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
	case core.WarehouseMAIN:
		return proto.NodeType_MAIN
	case core.WarehouseCENTER:
		return proto.NodeType_CENTRAL
	case core.WarehouseFREE:
		return proto.NodeType_FREE
	case core.WarehouseMALL:
		return proto.NodeType_MALL
	default:
		return proto.NodeType_UNRECOGNIZED
	}
}
