package execution

import (
	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/pkg/vector"
)

type engine struct {
	sorage       *IStorage
	space        *ISpace
	dimentions   uint
	learningRate float64
}

// Create new recommendation engine with configuration
func NewEngine(dimensions uint) IEngine {
	e := new(engine)
	spaceV1 := NewSpace(dimensions)
	storageV1 := NewStorage()

	e.sorage = &storageV1
	e.space = &spaceV1
	e.dimentions = dimensions
	e.learningRate = 0.01

	return e
}

// Add node to the recommendation system and get its uuid
func (e *engine) AddNode() uuid.UUID {
	nodeUuid := uuid.New()

	spaceNode := (*e.space).AddNode(nodeUuid)
	(*e.sorage).AddNode(nodeUuid, spaceNode)

	return nodeUuid
}

// Remove node by its uuid
func (e *engine) RemoveNode(nodeUuid uuid.UUID) error {
	spaceNode, err := (*e.sorage).RemoveNode(nodeUuid)
	if err != nil {
		return err
	}

	(*e.space).RemoveNode(spaceNode)

	return err
}

// Return map of nodes connections {<uuid>: <connection-weight>}
func (e *engine) GetNodeConnections(nodeUuid uuid.UUID) (map[uuid.UUID]uint, error) {
	return (*e.sorage).GetNodeConnections(nodeUuid)
}

// Add connection between nodes by thire uuids
func (e *engine) AddConnection(from uuid.UUID, to uuid.UUID, weight uint) error {
	return (*e.sorage).AddConnection(from, to, weight)
}

// Remove connection between nodes by thire uuids
func (e *engine) RemoveConnection(from uuid.UUID, to uuid.UUID) error {
	return (*e.sorage).RemoveConnection(from, to)
}

// Update connection weight
func (e *engine) UpdateConnectionWeight(from uuid.UUID, to uuid.UUID, weight uint) error {
	return (*e.sorage).UpdateConnectionWeight(from, to, weight)
}

// Optimaize recommendations by updating the vector space
func (e *engine) Optimize() {
	nodeUuids := (*e.sorage).GetNodesUUIDChannel()

	for nodeUuid := range nodeUuids {
		spaceNode, _ := (*e.sorage).GetSpaceNode(nodeUuid)
		currentPosition := (*spaceNode).GetVector()
		connections, _ := (*e.sorage).GetNodeConnections(nodeUuid)
		var updatedPosition vector.Vector

		for connectionNodeUuid, weight := range connections {
			connectionSpaceNode, _ := (*e.sorage).GetSpaceNode(connectionNodeUuid)
			diffVec := (*connectionSpaceNode).GetVector().Sub(currentPosition)
			diffNorm := diffVec.Norm()
			normalizedDiffVec := diffVec.Normalize()
			shiftLenght := (diffNorm - float64(weight)) * e.learningRate
			updatedPosition = currentPosition.Add(normalizedDiffVec.Mul(shiftLenght))
		}

		updatedSpaceNode := NewSpaceNodeFromVector(updatedPosition, nodeUuid)
		(*e.sorage).UpdateSpaceNode(nodeUuid, &updatedSpaceNode)
		(*e.space).UpdateNode(spaceNode, &updatedSpaceNode)
	}
}

// Get k best recommendations for node
func (e *engine) Recommend(nodeUuid uuid.UUID, k int) ([]uuid.UUID, error) {
	spaceNode, err := (*e.sorage).GetSpaceNode(nodeUuid)
	if err != nil {
		return nil, err
	}

	return (*e.space).KNN(spaceNode, k), nil
}
