package dag

import (
	"container/list"
	"errors"
	"log"
)

type Vertex struct {
	Parents  *list.List
	Children *list.List
	Value    interface{}
	Hash     string
	Type     string
}

func NewVertex(vType, hash string, value interface{}) *Vertex {
	return &Vertex{
		Hash:     hash,
		Type:     vType,
		Value:    value,
		Parents:  list.New(),
		Children: list.New(),
	}
}

func (v *Vertex) RemoveChild(hash string) {
	for e := v.Children.Front(); e != nil; e = e.Next() {
		val, ok := e.Value.(*Vertex)
		if !ok {
			log.Println(errors.New("type error"))
		}

		if val.Hash == hash {
			v.Children.Remove(e)
		}

	}
}

func (v *Vertex) RemoveParent(hash string) {
	for e := v.Parents.Front(); e != nil; e = e.Next() {
		val, ok := e.Value.(*Vertex)
		if !ok {
			log.Println(errors.New("type error"))
		}

		if val.Hash == hash {
			v.Parents.Remove(e)
		}

	}
}

func (v *Vertex) IsEqual(vC *Vertex) bool {
	if v.Hash != vC.Hash {
		return false
	}

	if v.Parents.Len() != vC.Parents.Len() {
		return false
	}

	e := v.Parents.Front()
	eC := vC.Parents.Front()
	for e != nil && eC != nil {

		if e.Value.(*Vertex).Hash != eC.Value.(*Vertex).Hash {
			return false
		}

		e.Next()
		eC.Next()
	}

	if v.Children.Len() != vC.Children.Len() {
		return false
	}

	return true
}
