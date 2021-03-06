package dag

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var (
	ErrCycle           = errors.New("dag: cycle between edges")
	ErrEdgeExists      = errors.New("dag: edge already exists")
	ErrVertexExists    = errors.New("dag: vertex already exists")
	ErrVertexNotExists = errors.New("dag: vertex does not exist")
)

type DAG struct {
	Vertexes *sync.Map //map[string]*Vertex
	Length   int
}

func NewDAG() *DAG {
	return &DAG{
		Vertexes: &sync.Map{}, // make(map[string]*Vertex),
		Length:   0,
	}
}

func (dag *DAG) AddVertex(vertex *Vertex) error {
	if _, ok := dag.Vertexes.Load(vertex.Hash); ok {
		return ErrVertexExists
	}

	dag.Vertexes.Store(vertex.Hash, vertex)
	dag.Length++

	return nil
}

func (dag *DAG) RemoveVertex(vertex *Vertex) {
	vItem, ok := dag.Vertexes.Load(vertex.Hash)
	if !ok {
		return
	}

	v := vItem.(*Vertex)

	for eParent := v.Parents.Front(); eParent != nil; eParent = eParent.Next() {
		eParent.Value.(*Vertex).RemoveChild(vertex.Hash)
	}

	for eChild := v.Children.Front(); eChild != nil; eChild = eChild.Next() {
		eChild.Value.(*Vertex).RemoveParent(vertex.Hash)
	}

	dag.Vertexes.Delete(vertex.Hash)
	dag.Length--
}

func (dag *DAG) AddEdge(from, to *Vertex) error {
	if from.Hash == to.Hash {
		return ErrCycle
	}

	fromVItem, ok := dag.Vertexes.Load(from.Hash)
	if !ok {
		return ErrVertexNotExists
	}
	fromV := fromVItem.(*Vertex)

	toVItem, ok := dag.Vertexes.Load(to.Hash)
	if !ok {
		return ErrVertexNotExists
	}
	toV := toVItem.(*Vertex)

	for e := fromV.Children.Front(); e != nil; e = e.Next() {

		//fmt.Printf("fromt pointer: %p %s, to pointer: %p %s, child pointer: %p %s \n", fromV, fromV.Hash, toV, toV.Hash, e.Value, e.Value.(*Vertex).Hash)

		if e.Value.(*Vertex).IsEqualPointer(fmt.Sprintf("%p", to)) {
			return ErrEdgeExists
		}
	}

	if dag.DepthFirstSearch(toV.Hash, fromV.Hash) {
		return ErrCycle
	}

	fromV.Children.PushBack(toV)
	toV.Parents.PushBack(fromV)

	return nil
}

func (dag *DAG) RemoveEdge(from, to *Vertex) error {

	fromVItem, ok := dag.Vertexes.Load(from.Hash)
	if !ok {
		return ErrVertexNotExists
	}
	fromV := fromVItem.(*Vertex)

	toVItem, ok := dag.Vertexes.Load(to.Hash)
	if !ok {
		return ErrVertexNotExists
	}
	toV := toVItem.(*Vertex)

	fromV.RemoveChild(toV.Hash)
	toV.RemoveParent(fromV.Hash)
	return nil
}

func (dag *DAG) EdgeExists(from, to *Vertex) (bool, error) {

	fromVItem, ok := dag.Vertexes.Load(from.Hash)
	if !ok {
		return false, ErrVertexNotExists
	}
	fromV := fromVItem.(*Vertex)

	toVItem, ok := dag.Vertexes.Load(to.Hash)
	if !ok {
		return false, ErrVertexNotExists
	}
	toV := toVItem.(*Vertex)

	// quick return
	if toV.Parents.Len() == 0 {
		return false, nil
	}

	for eChild := fromV.Children.Front(); eChild != nil; eChild = eChild.Next() {
		if eChild.Value.(*Vertex).IsEqualPointer(fmt.Sprintf("%p", toV)) {
			return true, nil
		}
	}

	return false, nil
}

func (dag *DAG) GetVertex(hash string) *Vertex {
	if v, ok := dag.Vertexes.Load(hash); ok {
		return v.(*Vertex)
	}

	return nil
}

func (dag *DAG) DepthFirstSearch(fromVertexHash, toVertexHash string) bool {
	found := map[string]bool{}
	dag.dfs(found, fromVertexHash)
	return found[toVertexHash]
}

func (dag *DAG) dfs(found map[string]bool, vertexId string) {
	vertexItem, ok := dag.Vertexes.Load(vertexId)
	if !ok {
		return
	}

	vertex := vertexItem.(*Vertex)
	for eChild := vertex.Children.Front(); eChild != nil; eChild = eChild.Next() {
		hash := eChild.Value.(*Vertex).Hash
		if !found[hash] {
			found[hash] = true
			dag.dfs(found, hash)
		}
	}

}

func (dag *DAG) IsEqual(dagC *DAG) (result bool) {

	if dag.Length != dagC.Length {
		return false
	}

	var check = func(vHash, value interface{}) bool {
		result = false

		v := value.(*Vertex)
		vCItem, ok := dagC.Vertexes.Load(vHash)
		if !ok {
			return result
		}
		vC := vCItem.(*Vertex)
		if !v.IsEqualPointer(fmt.Sprintf("%p", vC)) {
			return result
		}

		result = true

		return result
	}

	dag.Vertexes.Range(check)

	return result
}

// Copy shallow Copy
func (dag *DAG) Copy() *DAG {
	dagNew := NewDAG()

	// copy vertexes
	dag.Vertexes.Range(func(hash, value interface{}) bool {
		dagNew.Vertexes.Store(hash, value)
		dagNew.Length++
		return true
	})

	// copy edges
	dag.Vertexes.Range(func(hash, value interface{}) bool {
		v := value.(*Vertex)
		for eChild := v.Children.Front(); eChild != nil; eChild = eChild.Next() {
			err := dagNew.AddEdge(v, eChild.Value.(*Vertex))
			if err != nil {
				//panic(err)
				log.Println(err)
				return false
			}
		}

		return true
	})

	return dagNew
}

func (dag *DAG) Print(rootHash string) (str string) {

	/*
		dag.Vertexes.Range(func(hash, value interface{}) bool {
			v := value.(*Vertex)

			if v.Children.Len() == 0 {
				str = "\n" + str + dag.print(v, "") + "\n"
			}

			return true
		})

	*/

	value, ok := dag.Vertexes.Load(rootHash)
	if !ok {
		return ""
	}
	v := value.(*Vertex)
	str = "\n" + str + dag.print(v, "") + "\n"

	return str
}

func (dag *DAG) print(root *Vertex, prefix string) string {

	str := prefix + root.Hash + "\n"

	for eChild := root.Children.Front(); eChild != nil; eChild = eChild.Next() {
		child := eChild.Value.(*Vertex)

		if child.Children.Len() == 0 {
			str = str + dag.print(child, prefix+"    ")
		} else {
			str = str + dag.print(child, prefix+"    |")
		}

	}

	return str
}

// TopologicalSort get the vertexes without parents
func (dag *DAG) TopologicalSort() []*Vertex {
	copyV := dag.Copy()

	var sort = make([]*Vertex, 0)
	copyV.Vertexes.Range(func(hash, value interface{}) bool {
		v := value.(*Vertex)

		if v.Parents.Len() != 0 {
			return true
		}
		for eChild := v.Children.Front(); eChild != nil; eChild = eChild.Next() {
			child := eChild.Value.(*Vertex)
			child.RemoveChild(v.Hash)
		}
		copyV.Vertexes.Delete(v.Hash)

		// get the vertex without parents, the first one is the ROOT Vertex
		sort = append(sort, v)
		return true
	})

	return sort
}

/*
func (dag *DAG) TopologicalSortStable() []*Vertex {
	copyV := dag.Copy()
	noParentsVertexes := NewSortedVertexes()
	length := len(copyV.Vertexes)
	sort := make([]*Vertex, 0, length)
	if length == 0 {
		return sort
	}

	for {
		for _, v := range copyV.Vertexes {
			if v.Parents.Len() != 0 {
				continue
			}
			// get the root vertex
			noParentsVertexes.Add(v)
			delete(copyV.Vertexes, v.Hash)
		}
		firstNoParentsVertex := noParentsVertexes.PopFront()
		sort = append(sort, firstNoParentsVertex)
		if len(sort) == length {
			break
		}
		for eChild := firstNoParentsVertex.Children.Front(); eChild != nil; eChild = eChild.Next() {
			child := eChild.Value.(*Vertex)
			child.RemoveChild(firstNoParentsVertex.Hash)
		}
	}

	return sort
}
*/
