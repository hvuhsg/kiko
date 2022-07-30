package execution

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type storage struct {
	nodes       map[uuid.UUID]*ISpaceNode
	connections map[uuid.UUID]map[uuid.UUID]uint
	lock        sync.Mutex
}

// Create new RAM storage for nodes and connections
func NewStorage() IStorage {
	s := new(storage)
	s.nodes = make(map[uuid.UUID]*ISpaceNode, 1000)
	s.connections = make(map[uuid.UUID]map[uuid.UUID]uint, 1000)
	return s
}

// Save node in the RAM
func (s *storage) AddNode(nodeUuid uuid.UUID, spaceNode *ISpaceNode) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.nodes[nodeUuid] = spaceNode
	s.connections[nodeUuid] = make(map[uuid.UUID]uint, 10)
}

// Remove node from the RAM storage
func (s *storage) RemoveNode(nodeUuid uuid.UUID) (*ISpaceNode, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	spaceNode, ok := s.nodes[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' not found", nodeUuid.String())
	}

	delete(s.nodes, nodeUuid)
	delete(s.connections, nodeUuid)

	return spaceNode, nil
}

// Update nodes spaceNode
func (s *storage) UpdateSpaceNode(nodeUuid uuid.UUID, spaceNode *ISpaceNode) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.nodes[nodeUuid] = spaceNode
}

func (s *storage) AddConnection(from uuid.UUID, to uuid.UUID, weight uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, fromExist := s.nodes[from]
	_, toExist := s.nodes[to]

	if !fromExist || !toExist {
		return fmt.Errorf("connection nodes does not exist")
	}

	_, ok := s.connections[from][to]
	if ok {
		return fmt.Errorf("connection already exist")
	}

	s.connections[from][to] = weight

	return nil
}

func (s *storage) RemoveConnection(from uuid.UUID, to uuid.UUID) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.connections[from][to]
	if !ok {
		return fmt.Errorf("connection does not exist")
	}

	delete(s.connections[from], to)

	return nil
}

func (s *storage) UpdateConnectionWeight(from uuid.UUID, to uuid.UUID, updatedWeight uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.connections[from][to]
	if !ok {
		return fmt.Errorf("connection does not exist")
	}

	s.connections[from][to] = updatedWeight

	return nil
}

func (s *storage) GetNodeConnections(nodeUuid uuid.UUID) (map[uuid.UUID]uint, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	connections, ok := s.connections[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' does not exist", nodeUuid.String())
	}

	return connections, nil
}

// Get node by uuid
func (s *storage) GetSpaceNode(nodeUuid uuid.UUID) (*ISpaceNode, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	spaceNode, ok := s.nodes[nodeUuid]
	if !ok {
		return nil, fmt.Errorf("node '%s' does not exist", nodeUuid.String())
	}

	return spaceNode, nil
}

func (s *storage) GetNodesUUIDChannel() chan uuid.UUID {
	uuidsChan := make(chan uuid.UUID, 50)

	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		defer close(uuidsChan)
		for nodeUuid := range s.nodes {
			uuidsChan <- nodeUuid
		}
	}()

	return uuidsChan
}
