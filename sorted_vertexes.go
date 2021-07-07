package dag

import "container/list"

type SortedVertexes struct {
	*list.List
}

func NewSortedVertexes() *SortedVertexes {
	l := list.New()
	return &SortedVertexes{l}
}

func (s *SortedVertexes) Add(v *Vertex) {
	for e := s.Front(); e != nil; e = e.Next() {
		if v.Hash < e.Value.(*Vertex).Hash {
			s.InsertBefore(v, e)
			return
		}
	}
	s.PushBack(v)
}

func (s *SortedVertexes) PopFront() *Vertex {
	e := s.Front()
	if nil == e {
		return nil
	}
	s.Remove(e)
	return e.Value.(*Vertex)
}
