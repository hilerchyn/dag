package dag

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewDAG(t *testing.T) {

	MAX := 100000000

	dag := NewDAG()
	root := &Vertex{Hash: "000000", Type: "genesis", Value: "the genesis block"}
	dag.AddVertex(root)
	for i := 0; i < MAX; i++ {
		v := &Vertex{Hash: fmt.Sprintf("%d", i), Type: "TX", Value: "transaction"}
		dag.AddVertex(v)
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

		dag.AddEdge(
			&Vertex{Hash: fmt.Sprintf("%d", from), Type: "TX", Value: "transaction"},
			&Vertex{Hash: fmt.Sprintf("%d", to), Type: "TX", Value: "transaction"},
		)
	}

	time.Sleep(10 * time.Second)

}
