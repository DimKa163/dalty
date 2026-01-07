package usecase

import (
	"context"
	"errors"
	"time"

	graph2 "github.com/DimKa163/dalty/pkg/graph"

	"github.com/DimKa163/dalty/internal/logging"
	"github.com/DimKa163/dalty/internal/warehouse/core"
	"github.com/beevik/guid"
	"go.uber.org/zap"
)

type PathService struct {
	warehouseRepository core.WarehouseRepository
	pathFinder          *core.PathFinder
	graphContext        *graph2.GraphContext
}

func NewPathService(warehouseRepository core.WarehouseRepository, pathFinder *core.PathFinder, graphContext *graph2.GraphContext) *PathService {
	return &PathService{warehouseRepository: warehouseRepository, pathFinder: pathFinder, graphContext: graphContext}
}

func (ps *PathService) GetPath(ctx context.Context, dest, defaultWh *guid.Guid) (*core.Path, error) {
	gr, err := ps.graphContext.Get(ctx)
	if err != nil {
		return nil, err
	}
	node, ok := gr.Find(dest.String())
	if !ok {
		return nil, errors.New("dest not found")
	}
	path, err := ps.pathFinder.Path(ctx, node)
	if err != nil {
		return nil, err
	}
	if !path.Contains(defaultWh.String()) {
		node, _ = gr.Find(defaultWh.String())
		path, err = ps.pathFinder.Path(ctx, node)
		if err != nil {
			return nil, err
		}
	}
	return path, nil
}

func (ps *PathService) UpdateGraph(ctx context.Context) error {
	logger := logging.Logger(ctx)
	logger.Info("start to update warehouse graph")
	gr := graph2.NewGraph()
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
		gr.AddNode(node)

		if w.SenderID != nil {
			wS, ok := warehouseMap[w.SenderID.String()]
			if !ok {
				logger.Warn("sender not found", zap.String("sender_id", w.SenderID.String()), zap.String("node", w.Name))
				continue
			}
			sender := createNode(wS)
			gr.AddNode(sender)
			gr.AddEdge(sender, node, 0)
			loggerSug.Debugf("%s send to %s", wS.Name, w.Name)
		}
		if w.RecipientID != nil {
			wR, ok := warehouseMap[w.RecipientID.String()]
			if !ok {
				logger.Warn("recipient not found", zap.String("recipient_id", w.RecipientID.String()), zap.String("node", w.Name))
				continue
			}
			recipient := createNode(wR)
			gr.AddNode(recipient)
			gr.AddEdge(node, recipient, 0)
			loggerSug.Debugf("%s send to %s", w.Name, wR.Name)
		}
	}
	ps.graphContext.Update(gr)
	elapsed := time.Since(startTime)
	logger.Info("warehouse graph updated successfully", zap.Duration("elapsed", elapsed))
	return nil
}

func createNode(w *core.Warehouse) *graph2.Node {
	var node graph2.Node
	node.ID = w.ID.String()
	node.Value = w
	return &node
}
