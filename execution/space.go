package execution

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/pkg/vector"
	"github.com/kyroy/kdtree"
)

type space struct {
	kdtreeSpace kdtree.KDTree
	dimensions  uint
}

// Create nodes vector space
func NewSpace(dimensions uint) ISpace {
	s := new(space)
	s.kdtreeSpace = *kdtree.New([]kdtree.Point{})
	s.dimensions = dimensions

	return s
}

// Add node to vector space
func (s *space) AddNode(nodeUuid uuid.UUID) *ISpaceNode {
	spaceNode := NewRandomSpaceNode(s.dimensions, nodeUuid)

	s.kdtreeSpace.Insert(spaceNode)

	return &spaceNode
}

// Remove node from vector space
func (s *space) RemoveNode(spaceNode *ISpaceNode) error {
	foundNode := s.kdtreeSpace.Remove(*spaceNode)
	if foundNode == nil {
		return fmt.Errorf("node '%s' does not exist", (*spaceNode).GetUUID().String())
	}

	return nil
}

func (s *space) UpdateNode(old *ISpaceNode, new *ISpaceNode) {
	s.kdtreeSpace.Remove(*old)
	s.kdtreeSpace.Insert(*new)
}

func (s *space) KNN(spaceNode *ISpaceNode, k int) []uuid.UUID {
	closeNodes := s.kdtreeSpace.KNN(*spaceNode, int(k+1))
	nodesUuids := make([]uuid.UUID, len(closeNodes)-1)

	// TODO: replace this stupid solution
	for index, node := range closeNodes {
		if index == 0 {
			continue
		}

		nodeUuid, _ := uuid.Parse(fmt.Sprint(node))
		nodesUuids[index-1] = nodeUuid
	}

	return nodesUuids
}

type spaceNode struct {
	vector vector.Vector
	uuid   uuid.UUID
}

// Create new node from vector and uuid
func NewSpaceNodeFromVector(vec vector.Vector, nodeUuid uuid.UUID) ISpaceNode {
	s := new(spaceNode)
	s.vector = vec
	s.uuid = nodeUuid

	return s
}

// Create new node from array and uuid
func NewSpaceNodeFromArray(arr []float64, nodeUuid uuid.UUID) ISpaceNode {
	s := new(spaceNode)
	s.vector = vector.NewVectorFromArray(arr)
	s.uuid = nodeUuid

	return s
}

// Create node with random vector with values in range 0-1
func NewRandomSpaceNode(lenght uint, nodeUuid uuid.UUID) ISpaceNode {
	s := new(spaceNode)
	s.vector = vector.NewRandomVector(int(lenght))
	s.uuid = nodeUuid

	return s
}

// Get number of node vector dimensions (vector lenght)
func (s *spaceNode) Dimensions() int {
	return s.vector.Dimensions()
}

// Get vector dimension value
func (s *spaceNode) Dimension(i int) float64 {
	return s.vector.Dimension(i)
}

// Get node string representaion (the node uuid as string)
func (s *spaceNode) String() string {
	return s.GetUUID().String()
}

// Get node uuid
func (s *spaceNode) GetUUID() uuid.UUID {
	return s.uuid
}

// Get node vector
func (s *spaceNode) GetVector() vector.Vector {
	return s.vector
}
