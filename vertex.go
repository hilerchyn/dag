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

	/*
		经过测试，无法使用指针地址类比对是否相等，原因如下：
		1、所有的顶点都存储与sync.Map中
		2、比对时使用的是顶点的Parents和Children类型为*list.List
		3、从sync.Map中Load顶点数据时，应该是值传递，或者 向 *list.List PushBack 时 使用的是值传递，导致无法指针匹配
		srcPointer := fmt.Sprintf("%p", v)
		destPointer := fmt.Sprintf("%p", vC)
		if strings.Compare(srcPointer, destPointer) == 0 {
			log.Println("compare pointer: ", srcPointer, destPointer, strings.Compare(srcPointer, destPointer))
		}
		if strings.Compare(v.Hash, vC.Hash) == 0 {
			log.Println("compare hash: ", v.Hash, vC.Hash, srcPointer, destPointer, strings.Compare(srcPointer, destPointer))
		}

		if strings.Compare(srcPointer, destPointer) != 0{
			return false
		}
		log.Println(srcPointer, destPointer)
	*/

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
