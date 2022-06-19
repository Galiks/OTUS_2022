package list

import (
	"fmt"
	"strings"
)

type List interface {
	Len() int
	Front() *Item
	Back() *Item
	PushFront(v interface{}) *Item
	PushBack(v interface{}) *Item
	Remove(i *Item)
	MoveToFront(i *Item)
	Print() string
}

type Item struct {
	Value interface{}
	Next  *Item
	Prev  *Item
}

type list struct {
	len   int
	first *Item
	last  *Item
}

func NewList() List {
	return &list{}
}

func (l *list) Print() string {
	var (
		elem   *Item
		result = new(strings.Builder)
	)
	for i := 0; i < l.Len(); i++ {
		if i == 0 {
			elem = l.first
			result.WriteString(fmt.Sprint(elem.Value))
		} else {
			result.WriteString(fmt.Sprint(elem.Next.Value))
			elem = elem.Next
		}
	}
	return result.String()
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *Item {
	return l.first
}

func (l *list) Back() *Item {
	return l.last
}

func (l *list) PushFront(v interface{}) *Item {
	newItem := &Item{
		Value: v,
		Prev:  nil,
	}
	exFirst := l.first
	l.first = newItem
	if l.Len() == 0 {
		l.last = l.first
	} else {
		exFirst.Prev = newItem
		newItem.Next = exFirst
	}

	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *Item {
	newItem := &Item{
		Value: v,
		Next:  nil,
	}
	exLast := l.last
	l.last = newItem
	if l.Len() == 0 {
		l.first = l.last
	} else {
		exLast.Next = newItem
		newItem.Prev = exLast
	}

	l.len++

	return newItem
}

func (l *list) Remove(i *Item) {
	l.removeLink(i)
}

func (l *list) MoveToFront(i *Item) {
	l.PushFront(i.Value)
	l.removeLink(i)
}

func (l *list) removeLink(i *Item) {
	next := i.Next
	prev := i.Prev
	if prev == nil {
		l.first = next
	} else {
		prev.Next = next
	}

	if next == nil {
		l.last = prev
	} else {
		next.Prev = prev
	}
	l.len--
}
