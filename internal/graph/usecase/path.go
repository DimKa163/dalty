package usecase

import (
	"context"
	"errors"
	"github.com/DimKa163/dalty/internal/graph/core"
	"github.com/DimKa163/dalty/internal/logging"
	"github.com/beevik/guid"
	"go.uber.org/zap"
	"time"
)

type PathService struct {
	warehouseRepository core.WarehouseRepository
	graphContext        *core.GraphContext
}

func NewPathService(warehouseRepository core.WarehouseRepository, graphContext *core.GraphContext) *PathService {
	return &PathService{warehouseRepository: warehouseRepository, graphContext: graphContext}
}

func (ps *PathService) GetPath(ctx context.Context, dest, defaultWh *guid.Guid) (*core.Path, error) {
	graph, err := ps.graphContext.Get(ctx)
	if err != nil {
		return nil, err
	}
	node, ok := graph.Find(dest)
	if !ok {
		return nil, errors.New("dest not found")
	}
	path := graph.Path(node)
	if !path.Contains(defaultWh) {
		node, _ = graph.Find(defaultWh)
		path = graph.Path(node)
	}
	return path, nil
}

func (ps *PathService) UpdateGraph(ctx context.Context) error {
	logger := logging.Logger(ctx)
	logger.Info("start to update graph")
	graph := core.NewGraph()
	startTime := time.Now()
	warehouses, err := ps.warehouseRepository.GetAll(ctx)
	if err != nil {
		logger.Error("error occurred when GetAll Warehouses", zap.Error(err))
		return err
	}
	warehouseMap := make(map[string]*core.Warehouse)
	for _, w := range warehouses {
		warehouseMap[w.ID.String()] = w
	}
	loggerSug := logger.Sugar()
	for _, w := range warehouses {
		node := createNode(w)
		graph.AddNode(node)

		if w.SenderID != nil {
			wS, ok := warehouseMap[w.SenderID.String()]
			if !ok {
				logger.Warn("sender not found", zap.String("sender_id", w.SenderID.String()), zap.String("node", node.Name))
				continue
			}
			sender := createNode(wS)
			graph.AddNode(sender)
			graph.AddEdge(sender, node, 0)
			loggerSug.Debugf("%s send to %s", sender.Name, node.Name)
		}
		if w.RecipientID != nil {
			wR, ok := warehouseMap[w.RecipientID.String()]
			if !ok {
				logger.Warn("recipient not found", zap.String("recipient_id", w.RecipientID.String()), zap.String("node", node.Name))
				continue
			}
			recipient := createNode(wR)
			graph.AddNode(recipient)
			graph.AddEdge(node, recipient, 0)
			loggerSug.Debugf("%s send to %s", node.Name, recipient.Name)
		}
	}
	ps.graphContext.Update(graph)
	elapsed := time.Since(startTime)
	logger.Info("graph updated successfully", zap.Duration("elapsed", elapsed))
	return nil
}

func createNode(w *core.Warehouse) *core.Node {
	var node core.Node
	node.ID = w.ID
	node.Name = w.Name
	node.OnlyStockPickupAllowed = w.OnlyStockPickupAllowed
	if w.Category != nil {
		node.AvailableRest = w.Category.AvailableForBalance
		node.Type = core.MapWarehouseType(w.Category.Fnrec)
	}
	if w.Info != nil {
		node.Address = w.Info.Address
		node.DescriptorGroup = w.Info.DescriptorGroup
		if w.Info.TimeZone != nil {
			node.Code = w.Info.TimeZone.Code
		}
	}
	return &node
}
