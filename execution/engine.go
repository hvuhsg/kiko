package execution

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/pkg/vector"
)

type engine struct {
	sorage                   *IStorage
	space                    *ISpace
	dimentions               uint
	learningRate             float64
	minOptimizationRoundTime time.Duration
	locker                   sync.Mutex
}

// Create new recommendation engine with configuration
func NewEngine(dimensions uint, learningRate float64, optimizationRoundDuration time.Duration) IEngine {
	e := new(engine)
	spaceV1 := NewSpace(dimensions)
	storageV1 := NewStorage()

	e.sorage = &storageV1
	e.space = &spaceV1
	e.dimentions = dimensions
	e.learningRate = learningRate
	e.minOptimizationRoundTime = optimizationRoundDuration

	return e
}

// Add node to the recommendation system and get its uuid
func (e *engine) AddNode() uuid.UUID {
	e.lock()
	defer e.unlock()

	nodeUuid := uuid.New()

	spaceNode := (*e.space).AddNode(nodeUuid)
	(*e.sorage).AddNode(nodeUuid, spaceNode)

	return nodeUuid
}

// Remove node by its uuid
func (e *engine) RemoveNode(nodeUuid uuid.UUID) error {
	e.lock()
	defer e.unlock()

	spaceNode, err := (*e.sorage).RemoveNode(nodeUuid)
	if err != nil {
		return err
	}

	(*e.space).RemoveNode(spaceNode)

	return err
}

// Return map of nodes connections {<uuid>: <connection-weight>}
func (e *engine) GetNodeConnections(nodeUuid uuid.UUID) (map[uuid.UUID]uint, error) {
	e.lock()
	defer e.unlock()

	return (*e.sorage).GetNodeConnections(nodeUuid)
}

// Add connection between nodes by thire uuids
func (e *engine) AddConnection(from uuid.UUID, to uuid.UUID, weight uint) error {
	e.lock()
	defer e.unlock()

	err := (*e.sorage).AddConnection(from, to, weight)
	if err != nil {
		return err
	}

	return (*e.sorage).AddConnection(to, from, weight)
}

// Remove connection between nodes by thire uuids
func (e *engine) RemoveConnection(from uuid.UUID, to uuid.UUID) error {
	e.lock()
	defer e.unlock()

	err := (*e.sorage).RemoveConnection(from, to)
	if err != nil {
		return err
	}

	return (*e.sorage).RemoveConnection(to, from)
}

// Update connection weight
func (e *engine) UpdateConnectionWeight(from uuid.UUID, to uuid.UUID, weight uint) error {
	e.lock()
	defer e.unlock()

	err := (*e.sorage).UpdateConnectionWeight(from, to, weight)
	if err != nil {
		return err
	}

	return (*e.sorage).UpdateConnectionWeight(to, from, weight)
}

func optimaizeNode(e *engine, nodeUuid uuid.UUID) {
	spaceNode, _ := (*e.sorage).GetSpaceNode(nodeUuid)
	currentPosition := (*spaceNode).GetVector()
	connections, _ := (*e.sorage).GetNodeConnections(nodeUuid)
	var updatedPosition vector.Vector

	for connectionNodeUuid, weight := range connections {
		connectionSpaceNode, _ := (*e.sorage).GetSpaceNode(connectionNodeUuid)
		diffVec := (*connectionSpaceNode).GetVector().Sub(currentPosition)
		diffNorm := diffVec.Norm()

		if diffNorm < float64(weight) {
			diffVec = diffVec.Mul(-1)
		}

		updatedPosition = currentPosition.Add(diffVec.Mul(e.learningRate))
	}

	if len(connections) > 0 {
		updatedSpaceNode := NewSpaceNodeFromVector(updatedPosition, nodeUuid)
		(*e.sorage).UpdateSpaceNode(nodeUuid, &updatedSpaceNode)
		(*e.space).UpdateNode(spaceNode, &updatedSpaceNode)
	}
}

func optimizationRound(e *engine) {
	e.lock()
	defer e.unlock()

	nodeUuids := (*e.sorage).GetNodes()

	for _, nodeUuid := range nodeUuids {
		optimaizeNode(e, nodeUuid)
	}
}

// Optimaize recommendations by updating the vector space
func (e *engine) Optimize() {
	for {
		optimizationRound(e)
		time.Sleep(e.minOptimizationRoundTime)
	}
}

// Get k best recommendations for node
func (e *engine) Recommend(nodeUuid uuid.UUID, k int) ([]uuid.UUID, error) {
	e.lock()
	defer e.unlock()

	spaceNode, err := (*e.sorage).GetSpaceNode(nodeUuid)
	if err != nil {
		return nil, err
	}

	return (*e.space).KNN(spaceNode, k), nil
}

func (e *engine) lock() {
	e.locker.Lock()
}

func (e *engine) unlock() {
	e.locker.Unlock()
}
