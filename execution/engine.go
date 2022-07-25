package execution

import (
	"github.com/google/uuid"
)

type engine struct {
	sorage     *IStorage
	space      *ISpace
	dimentions uint
}

func NewEngine(dimensions uint) IEngine {
	e := new(engine)
	spaceV1 := NewSpace(dimensions)
	storageV1 := NewStorage()

	e.sorage = &storageV1
	e.space = &spaceV1
	e.dimentions = dimensions

	return e
}

func (e engine) AddNode() uuid.UUID {
	nodeUuid := uuid.New()

	spaceNode := (*e.space).AddNode(nodeUuid)
	(*e.sorage).AddNode(nodeUuid, spaceNode)

	return nodeUuid
}

func (e engine) RemoveNode(nodeUuid uuid.UUID) error {
	spaceNode, err := (*e.sorage).RemoveNode(nodeUuid)
	if err != nil {
		return err
	}

	(*e.space).RemoveNode(spaceNode)

	return err
}

func (e engine) GetNodeConnections(nodeUuid uuid.UUID) (map[uuid.UUID]uint, error) {
	return (*e.sorage).GetNodeConnections(nodeUuid)
}

func (e engine) AddConnection(from uuid.UUID, to uuid.UUID, weight uint) error {
	return (*e.sorage).AddConnection(from, to, weight)
}

func (e engine) RemoveConnection(from uuid.UUID, to uuid.UUID) error {
	return (*e.sorage).RemoveConnection(from, to)
}

func (e engine) UpdateConnectionWeight(from uuid.UUID, to uuid.UUID, weight uint) error {
	return (*e.sorage).UpdateConnectionWeight(from, to, weight)
}

func (e engine) Optimize() {}

func (e engine) KNN(nodeUuid uuid.UUID, k int) ([]uuid.UUID, error) {
	spaceNode, err := (*e.sorage).GetSpaceNode(nodeUuid)
	if err != nil {
		return nil, err
	}

	return (*e.space).KNN(spaceNode, k), nil
}
