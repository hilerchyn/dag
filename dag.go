package dag

import "errors"

var (
	ErrCycle           = errors.New("dag: cycle between edges")
	ErrEdgeExists      = errors.New("dag: edge already exists")
	ErrVertexExists    = errors.New("dag: vertex already exists")
	ErrVertexNotExists = errors.New("dag: vertex does not exist")
)

type DAG struct {
	Vertexes map[string]*Vertex
}

func NewDAG() *DAG {
	return &DAG{
		Vertexes: make(map[string]*Vertex),
	}
}

func (dag *DAG) AddVertex(vertex *Vertex) error {
	if _, ok := dag.Vertexes[vertex.Hash]; ok {
		return ErrVertexExists
	}

	dag.Vertexes[vertex.Hash] = vertex

	return nil
}

func (dag *DAG) RemoveVertex(vertex *Vertex) {
	v, ok := dag.Vertexes[vertex.Hash]
	if !ok {
		return
	}

	for eParent := v.Parents.Front(); eParent != nil; eParent.Next() {
		eParent.Value.(*Vertex).RemoveChild(vertex.Hash)
	}

	for eChild := v.Children.Front(); eChild != nil; eChild.Next() {
		eChild.Value.(*Vertex).RemoveParent(vertex.Hash)
	}

	delete(dag.Vertexes, vertex.Hash)

}

func (dag *DAG) AddEdge(from, to *Vertex) error {
	if from.Hash == to.Hash {
		return ErrCycle
	}

	fromV, ok := dag.Vertexes[from.Hash]
	if !ok {
		return ErrVertexNotExists
	}

	toV, ok := dag.Vertexes[to.Hash]
	if !ok {
		return ErrVertexNotExists
	}

	for e := fromV.Children.Front(); e != nil; e.Next() {
		if e.Value.(*Vertex).IsEqual(to) {
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

	fromV, ok := dag.Vertexes[from.Hash]
	if !ok {
		return ErrVertexNotExists
	}

	toV, ok := dag.Vertexes[to.Hash]
	if !ok {
		return ErrVertexNotExists
	}

	fromV.RemoveChild(toV.Hash)
	toV.RemoveParent(fromV.Hash)
	return nil
}

func (dag *DAG) EdgeExists(from, to *Vertex) (bool, error) {

	fromV, ok := dag.Vertexes[from.Hash]
	if !ok {
		return false, ErrVertexNotExists
	}

	toV, ok := dag.Vertexes[to.Hash]
	if !ok {
		return false, ErrVertexNotExists
	}

	// quick return
	if toV.Parents.Len() == 0 {
		return false, nil
	}

	for eChild := fromV.Children.Front(); eChild != nil; eChild.Next() {
		if eChild.Value.(*Vertex).IsEqual(toV) {
			return true, nil
		}
	}

	return false, nil
}

func (dag *DAG) GetVertex(hash string) *Vertex {
	if v, ok := dag.Vertexes[hash]; ok {
		return v
	}

	return nil
}

func (dag *DAG) DepthFirstSearch(fromVertexHash, toVertexHash string) bool {
	found := map[string]bool{}
	dag.dfs(found, fromVertexHash)
	return found[toVertexHash]
}

func (dag *DAG) dfs(found map[string]bool, vertexId string) {
	vertex, ok := dag.Vertexes[vertexId]
	if !ok {
		return
	}

	for eChild := vertex.Children.Front(); eChild != nil; eChild.Next() {
		hash := eChild.Value.(*Vertex).Hash
		if !found[hash] {
			found[hash] = true
			dag.dfs(found, hash)
		}
	}

}

func (dag *DAG) IsEqual(dagC *DAG) bool {
	if len(dag.Vertexes) != len(dagC.Vertexes) {
		return false
	}
	for vHash, v := range dag.Vertexes {
		vC, ok := dagC.Vertexes[vHash]
		if !ok {
			return false
		}
		if !v.IsEqual(vC) {
			return false
		}
	}
	return true
}

// Copy shallow Copy
func (dag *DAG) Copy() *DAG {
	dagNew := NewDAG()

	// copy vertexes
	for _, v := range dag.Vertexes {
		dagNew.Vertexes[v.Hash] = &Vertex{
			Hash:  v.Hash,
			Value: v.Value,
			Type:  v.Type,
		}
	}

	// copy edges
	for _, v := range dag.Vertexes {
		for eChild := v.Children.Front(); eChild != nil; eChild.Next() {
			err := dagNew.AddEdge(v, eChild.Value.(*Vertex))
			if err != nil {
				panic(err)
				//return nil
			}
		}
	}
	return dagNew
}

func (dag *DAG) Print() (str string) {
	for _, v := range dag.Vertexes {
		if v.Parents.Len() == 0 {
			str = str + dag.print(v, "") + "\n"
		}
	}
	return str
}

func (dag *DAG) print(root *Vertex, prefix string) string {
	str := prefix + root.Hash + "\n"
	for eChild := root.Children.Front(); ; {
		child := eChild.Value.(*Vertex)
		// the last element
		if eChild == nil {
			str = str + dag.print(child, prefix+"    ")
			break
		} else {
			str = str + dag.print(child, prefix+"    |")
		}

		eChild.Next()
	}
	return str
}

// TopologicalSort get the vertexes without parents
func (dag *DAG) TopologicalSort() []*Vertex {
	copyV := dag.Copy()

	var sort = make([]*Vertex, 0)
	for {
		for _, v := range copyV.Vertexes {
			if v.Parents.Len() != 0 {
				continue
			}
			for eChild := v.Children.Front(); eChild != nil; eChild.Next() {
				child := eChild.Value.(*Vertex)
				child.RemoveChild(v.Hash)
			}
			delete(copyV.Vertexes, v.Hash)

			// get the vertex without parents, the first one is the ROOT Vertex
			sort = append(sort, v)
		}
		if len(copyV.Vertexes) == 0 {
			break
		}
	}

	return sort
}

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
		for eChild := firstNoParentsVertex.Children.Front(); eChild != nil; eChild.Next() {
			child := eChild.Value.(*Vertex)
			child.RemoveChild(firstNoParentsVertex.Hash)
		}
	}

	return sort
}
