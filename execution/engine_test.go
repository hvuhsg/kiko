package execution_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/execution"
)

func TestAddNode(t *testing.T) {
	e := execution.NewEngine(2, 0.1, time.Minute*1)
	n1 := e.AddNode()

	_, err := e.GetNodeConnections(n1)

	if err != nil {
		t.Error(err)
	}
}

func TestRemoveExistingNode(t *testing.T) {
	e := execution.NewEngine(2, 0.1, time.Minute*1)
	n1 := e.AddNode()
	err := e.RemoveNode(n1)

	if err != nil {
		t.Errorf("error when trying to remove existing node")
	}
}

func TestRemoveNonExistingNode(t *testing.T) {
	e := execution.NewEngine(2, 0.1, time.Minute*1)
	err := e.RemoveNode(uuid.New())

	if err == nil {
		t.Errorf("must return error when trying to remove non existing node")
		return
	}
}

func TestOptimaize(t *testing.T) {
	e := execution.NewEngine(2, 0.1, time.Millisecond*25)
	n1 := e.AddNode()
	n2 := e.AddNode()
	n3 := e.AddNode()
	e.AddConnection(n1, n2, 1)
	e.AddConnection(n2, n1, 1)
	e.AddConnection(n1, n3, 10)
	e.AddConnection(n3, n1, 10)

	go e.Optimize()
	time.Sleep(1 * time.Second)

	results, err := e.Recommend(n1, 1)
	if err != nil {
		t.Error(err)
	}

	if results[0] != n2 {
		t.Error("must recommand n2")
	}
}
