package list

import (
	"Goddess/consts"
)

type node struct {
	value      interface{}
	prev, next *node
}

type LinkedList struct {
	first *node
	last  *node
	size  int
}

func (list *LinkedList) Push(where int, value interface{}) {
	node := &node{
		value: value,
		prev:  nil,
		next:  nil,
	}
	if where == consts.HEAD {
		node.next = list.first
		list.first.prev = node
		list.first = node
	} else if where == consts.TAIL {
		node.prev = list.last
		list.last.next = node
		list.last = node
	}
	list.size++
}

func (list *LinkedList) LPushCmd(value interface{}) {
	list.Push(consts.HEAD, value)
}

func (list *LinkedList) RPushCmd(value interface{}) {
	list.Push(consts.TAIL, value)
}
