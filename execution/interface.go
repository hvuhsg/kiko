package execution

import (
	"github.com/google/uuid"
	"github.com/kyroy/kdtree"
)

type ISpaceNode interface {
	kdtree.Point
	String() string
	GetUUID() uuid.UUID
	GetVector() []float64
}

// Presistence storage interface can use SQL / NoSQL DB
type IStorage interface {
	AddNode(uuid.UUID, *ISpaceNode)
	RemoveNode(uuid.UUID) (*ISpaceNode, error)
	AddConnection(uuid.UUID, uuid.UUID, uint) error
	RemoveConnection(uuid.UUID, uuid.UUID) error
	UpdateConnectionWeight(uuid.UUID, uuid.UUID, uint) error
	GetNodeConnections(uuid.UUID) (map[uuid.UUID]uint, error)
	GetSpaceNode(uuid.UUID) (*ISpaceNode, error)
}

type ISpace interface {
	AddNode(uuid.UUID) *ISpaceNode
	RemoveNode(*ISpaceNode) error

	// get k closest nodes to node sorted from closest to farsest
	KNN(*ISpaceNode, int) []uuid.UUID
}

type IEngine interface {
	AddNode() uuid.UUID
	RemoveNode(uuid.UUID) error
	GetNodeConnections(uuid.UUID) (map[uuid.UUID]uint, error)
	AddConnection(uuid.UUID, uuid.UUID, uint) error
	RemoveConnection(uuid.UUID, uuid.UUID) error
	UpdateConnectionWeight(uuid.UUID, uuid.UUID, uint) error

	Optimize()

	// get k most recommended nodes for node
	Recommend(uuid.UUID, int) ([]uuid.UUID, error)
}
