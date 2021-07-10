package dag

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"strings"
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

	/*
		因为go中没有指针传递，所以此处传入的vC的指针已经产生变化，因此无法使用指针比对
	*/

	fmt.Printf("child pointer: %s %p, to pointer: %s %p\n", v.Hash, v, vC.Hash, vC)

	// 因此此处，只通过比对Hash值来判定该值是否相等，
	if v.Hash == vC.Hash {
		return true
	}

	/*
		if v.Hash != vC.Hash {
			return true
		}
	*/

	if v.Parents.Len() != vC.Parents.Len() {
		return false
	}

	if v.Children.Len() != vC.Children.Len() {
		return false
	}

	e := v.Parents.Front()
	eC := vC.Parents.Front()
	for e != nil && eC != nil {

		if e.Value.(*Vertex).Hash != eC.Value.(*Vertex).Hash {
			return false
		}

		e = e.Next()
		eC = eC.Next()
	}

	return true
}

// IsEqualPointer 使用指针地址进行比较
func (v *Vertex) IsEqualPointer(comparePointer string) bool {

	if strings.Compare(fmt.Sprintf("%p", v), comparePointer) == 0 {
		return true
	}

	return false
}
