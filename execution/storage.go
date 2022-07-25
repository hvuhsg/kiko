package execution

import (
	"fmt"

	"github.com/google/uuid"
)

type storage struct {
	nodes       map[uuid.UUID]*ISpaceNode
	connections map[uuid.UUID]map[uuid.UUID]uint
}

// Ram nodes and connections storage
func NewStorage() IStorage {
	s := new(storage)
	s.nodes = make(map[uuid.UUID]*ISpaceNode, 1000)
	s.connections = make(map[uuid.UUID]map[uuid.UUID]uint, 1000)
	return s
}

func (s storage) AddNode(nodeUuid uuid.UUID, spaceNode *ISpaceNode) {
	s.nodes[nodeUuid] = spaceNode
	s.connections[nodeUuid] = make(map[uuid.UUID]uint, 10)
}

func (s storage) RemoveNode(nodeUuid uuid.UUID) (*ISpaceNode, error) {
	spaceNode, ok := s.nodes[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' not found", nodeUuid.String())
	}

	delete(s.nodes, nodeUuid)
	delete(s.connections, nodeUuid)

	return spaceNode, nil
}

func (s storage) AddConnection(from uuid.UUID, to uuid.UUID, weight uint) error {
	_, fromExist := s.nodes[from]
	_, toExist := s.nodes[to]

	if !fromExist || !toExist {
		return fmt.Errorf("connection nodes does not exist")
	}

	s.connections[from][to] = weight

	return nil
}

func (s storage) RemoveConnection(from uuid.UUID, to uuid.UUID) error {
	_, ok := s.connections[from][to]
	if !ok {
		return fmt.Errorf("connection does not exist")
	}

	delete(s.connections[from], to)

	return nil
}

func (s storage) UpdateConnectionWeight(from uuid.UUID, to uuid.UUID, updatedWeight uint) error {
	_, ok := s.connections[from][to]
	if !ok {
		return fmt.Errorf("connection does not exist")
	}

	s.connections[from][to] = updatedWeight

	return nil
}

func (s storage) GetNodeConnections(nodeUuid uuid.UUID) (map[uuid.UUID]uint, error) {
	connections, ok := s.connections[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' does not exist", nodeUuid.String())
	}

	return connections, nil
}

func (s storage) GetSpaceNode(nodeUuid uuid.UUID) (*ISpaceNode, error) {
	spaceNode, ok := s.nodes[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' does not exist", nodeUuid.String())
	}

	return spaceNode, nil
}
