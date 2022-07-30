package execution

import (
	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/pkg/vector"
	"github.com/kyroy/kdtree"
)

// Node uuid and vector data
type ISpaceNode interface {
	kdtree.Point
	String() string
	GetUUID() uuid.UUID
	GetVector() vector.Vector
}

// Presistence storage interface can use SQL / NoSQL DB
type IStorage interface {
	AddNode(uuid.UUID, *ISpaceNode)
	RemoveNode(uuid.UUID) (*ISpaceNode, error)
	UpdateSpaceNode(uuid.UUID, *ISpaceNode)
	AddConnection(uuid.UUID, uuid.UUID, uint) error
	RemoveConnection(uuid.UUID, uuid.UUID) error
	UpdateConnectionWeight(uuid.UUID, uuid.UUID, uint) error
	GetNodeConnections(uuid.UUID) (map[uuid.UUID]uint, error)
	GetSpaceNode(uuid.UUID) (*ISpaceNode, error)
	GetNodes() []uuid.UUID
}

// Vector space for the nodes
type ISpace interface {
	AddNode(uuid.UUID) *ISpaceNode
	RemoveNode(*ISpaceNode) error
	UpdateNode(*ISpaceNode, *ISpaceNode)

	// get k closest nodes to node sorted from closest to farsest
	KNN(*ISpaceNode, int) []uuid.UUID
}

// Recommendation system that works on a graph,
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
