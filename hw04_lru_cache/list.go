package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value      interface{}
	Next, Prev *ListItem
}

type list struct {
	len        int
	head, back *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++
	li := &ListItem{
		Value: v,
		Next:  l.Front(),
	}

	l.head = li

	if l.Back() == nil {
		l.back = li
	}

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	li := &ListItem{
		Value: v,
		Prev:  l.Back(),
	}

	if l.Back() != nil {
		l.Back().Next = li
	}

	l.back = li

	return li
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if l.head == i {
		l.head = i.Next
	}

	if l.back == i {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	if l.back == i {
		l.back = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.head.Prev = i
	i.Prev = nil
	i.Next = l.head
	l.head = i
}

func NewList() List {
	return &list{}
}
