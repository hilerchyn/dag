package dag

import (
	"fmt"
	"strings"
	"sync"
)

type Vertex struct {
	// @todo 使用sync.Map替换 list.List 提高监测速度
	Parents  *sync.Map
	Children *sync.Map
	Value    interface{}
	Hash     string
	Type     string
}

func NewVertex(vType, hash string, value interface{}) *Vertex {
	return &Vertex{
		Hash:     hash,
		Type:     vType,
		Value:    value,
		Parents:  new(sync.Map),
		Children: new(sync.Map),
	}
}

func (v *Vertex) RemoveChild(hash string) {

	v.Children.Delete(hash)
}

func (v *Vertex) RemoveParent(hash string) {
	v.Parents.Delete(hash)
}

// IsEqualPointer 使用指针地址进行比较
func (v *Vertex) IsEqualPointer(comparePointer string) bool {

	if strings.Compare(fmt.Sprintf("%p", v), comparePointer) == 0 {
		return true
	}

	return false
}
