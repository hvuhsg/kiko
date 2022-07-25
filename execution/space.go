package execution

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/kyroy/kdtree"
)

type space struct {
	kdtreeSpace kdtree.KDTree
	dimensions  uint
}

func NewSpace(dimensions uint) ISpace {
	s := new(space)
	s.kdtreeSpace = *kdtree.New([]kdtree.Point{})
	s.dimensions = dimensions

	return s
}

func (s space) AddNode(nodeUuid uuid.UUID) *ISpaceNode {
	spaceNode := NewRandomSpaceNode(s.dimensions, nodeUuid)

	s.kdtreeSpace.Insert(spaceNode)

	return &spaceNode
}

func (s space) RemoveNode(spaceNode *ISpaceNode) error {
	foundNode := s.kdtreeSpace.Remove(*spaceNode)
	if foundNode == nil {
		return fmt.Errorf("node '%s' does not exist", (*spaceNode).GetUUID().String())
	}

	return nil
}

func (s space) KNN(spaceNode *ISpaceNode, k int) []uuid.UUID {
	closeNodes := s.kdtreeSpace.KNN(*spaceNode, int(k))
	fmt.Println(closeNodes)

	nodesUuids := make([]uuid.UUID, len(closeNodes))

	return nodesUuids
}

type spaceNode struct {
	vector []float64
	uuid   uuid.UUID
}

func NewSpaceNode(vec []float64, nodeUuid uuid.UUID) ISpaceNode {
	s := new(spaceNode)
	s.vector = vec
	s.uuid = nodeUuid

	return s
}

func NewRandomSpaceNode(lenght uint, nodeUuid uuid.UUID) ISpaceNode {
	vec := make([]float64, lenght)

	for i := range vec {
		vec[i] = rand.Float64()
	}

	return NewSpaceNode(vec, nodeUuid)
}

func (s spaceNode) Dimensions() int {
	return len(s.vector)
}

func (s spaceNode) Dimension(i int) float64 {
	return s.vector[i]
}

func (s spaceNode) GetUUID() uuid.UUID {
	return s.uuid
}

func (s spaceNode) GetVector() []float64 {
	return s.vector
}
