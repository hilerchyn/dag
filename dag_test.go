package dag

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewDAG(t *testing.T) {

	MAX := 100000000

	dagNew := NewDAG()
	root := NewVertex("genesis", "000000", "the genesis block")
	dagNew.AddVertex(root)
	for i := 0; i < MAX; i++ {
		v := NewVertex("TX", fmt.Sprintf("%d", i), "transaction")
		dagNew.AddVertex(v)
	}

	for i := 0; i < MAX; i++ {
		from := rand.Intn(MAX)
		var to int
		for {
			to = rand.Intn(MAX)
			if to != from {
				break
			}
		}

		t.Log(from, to)

		dagNew.AddEdge(
			NewVertex("TX", fmt.Sprintf("%d", from), "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
			NewVertex("TX", fmt.Sprintf("%d", to), "transaction"),
		)
	}

	time.Sleep(10 * time.Second)

}

func TestDAG_IsEqual(t *testing.T) {
	MAX := 50

	dagNew := NewDAG()
	root := NewVertex("genesis", "000000", "the genesis block")
	dagNew.AddVertex(root)
	for i := 0; i < MAX; i++ {
		v := NewVertex("TX", fmt.Sprintf("%d", i), "transaction")
		dagNew.AddVertex(v)
	}

	for from := 0; from < MAX-1; from++ {
		for to := 1; to < MAX; to++ {
			dagNew.AddEdge(
				NewVertex("TX", fmt.Sprintf("%d", from), "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
				NewVertex("TX", fmt.Sprintf("%d", to), "transaction"),
			)
		}
	}

	dagNew.AddEdge(
		NewVertex("TX", fmt.Sprintf("%d", 10), "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
		NewVertex("TX", fmt.Sprintf("%d", 40), "transaction"),
	)

}

func TestNewSortedVertexes(t *testing.T) {

	MAX := 5

	dagNew := NewDAG()
	root := NewVertex("genesis", "000000", "the genesis block")
	dagNew.AddVertex(root)
	for i := 0; i < MAX; i++ {
		if i == 0 {
			continue
		}
		v := NewVertex("TX", fmt.Sprintf("%d", i), "transaction")
		dagNew.AddVertex(v)
	}

	for from := 0; from < MAX-1; from++ {
		fromHash := fmt.Sprintf("%d", from)
		if from == 0 {
			fromHash = "000000"
		}
		for to := 1; to < MAX; to++ {
			dagNew.AddEdge(
				NewVertex("TX", fromHash, "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
				NewVertex("TX", fmt.Sprintf("%d", to), "transaction"),
			)
		}
	}

	t.Log(dagNew.AddEdge(
		NewVertex("TX", fmt.Sprintf("%d", 2), "transaction"), // {Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
		NewVertex("TX", fmt.Sprintf("%d", 4), "transaction"),
	))

	t.Log(dagNew.Print("000000"))

}
