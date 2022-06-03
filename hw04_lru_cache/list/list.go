package list

import (
	"fmt"
	"strings"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Print() string
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func NewList() List {
	return &list{}
}

func (l *list) Print() string {
	var (
		elem   *ListItem
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

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
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

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
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

func (l *list) Remove(i *ListItem) {
	l.removeLink(i)
	l.len--
}

func (l *list) removeLink(i *ListItem) {
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
}

func (l *list) MoveToFront(i *ListItem) {
	l.removeLink(i)
	l.PushFront(i.Value)
}
