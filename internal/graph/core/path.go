package core

import (
	"container/list"
	"github.com/beevik/guid"
)

type Path struct {
	list  *list.List
	nodes map[guid.Guid]*list.Element
}

func NewPath() *Path {
	return &Path{
		list:  list.New(),
		nodes: make(map[guid.Guid]*list.Element),
	}
}

func (path *Path) Len() int {
	return path.list.Len()
}

func (path *Path) GetList() []*PathNode {

	result := make([]*PathNode, path.Len())
	var index int
	for e := path.list.Front(); e != nil; e = e.Next() {
		result[index] = e.Value.(*PathNode)
		index++
	}
	return result
}

func (path *Path) AddNode(n *PathNode) {
	el := path.list.PushBack(n)
	path.nodes[n.ID] = el
}

func (path *Path) AddAfter(n *PathNode, mark *PathNode) {
	elMark, ok := path.nodes[mark.ID]
	if !ok {
		return
	}
	path.list.InsertAfter(n, elMark)
}

func (path *Path) Contains(nodeId *guid.Guid) bool {
	_, ok := path.nodes[*nodeId]
	return ok
}

func (path *Path) FirstNode() *PathNode {
	if path.list.Len() == 0 {
		return nil
	}
	return path.list.Front().Value.(*PathNode)
}

func (path *Path) LastNode() *PathNode {
	if path.list.Len() == 0 {
		return nil
	}
	return path.list.Back().Value.(*PathNode)
}

type PathNode struct {
	Level int
	Next  *PathNode
	*Node
}

type DeliveryPath struct {
	First *PathNode
	Last  *PathNode
	Path  *Path
}

func NewDeliveryPath(path *Path) *DeliveryPath {
	return &DeliveryPath{
		First: path.FirstNode(),
		Last:  path.LastNode(),
		Path:  path,
	}
}

type DeliveryPathResult struct {
	Success bool
	Path    *DeliveryPath
}

func NewSuccessDeliveryPathResult(path *Path) *DeliveryPathResult {
	return &DeliveryPathResult{
		Success: true,
		Path:    NewDeliveryPath(path),
	}
}

func NewUnsuccessDeliveryPathResult(path *Path) *DeliveryPathResult {
	return &DeliveryPathResult{
		Success: false,
		Path:    NewDeliveryPath(path),
	}
}
