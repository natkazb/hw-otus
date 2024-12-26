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
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	start *ListItem
	end   *ListItem
	len   int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.start
}

func (l *list) Back() *ListItem {
	return l.end
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Next: l.start, Prev: nil}
	if l.start != nil {
		l.start.Prev = &newItem
	}
	l.start = &newItem
	if l.end == nil {
		l.end = l.start
	}
	l.len++
	return l.start
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Next: nil, Prev: l.end}
	if l.end != nil {
		l.end.Next = &newItem
	}
	l.end = &newItem
	if l.start == nil {
		l.start = l.end
	}
	l.len++
	return l.end
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.start = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.end = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.start {
		return
	}

	// Remove item from its current position
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.end = i.Prev
	}

	// Move to front
	i.Next = l.start
	i.Prev = nil
	l.start.Prev = i
	l.start = i
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.start
}

func (l *list) Back() *ListItem {
	return l.end
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Next: l.start, Prev: nil}
	if l.start != nil {
		l.start.Prev = &newItem
	}
	l.start = &newItem
	if l.end == nil {
		l.end = l.start
	}
	l.len++
	return l.start
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Next: nil, Prev: l.end}
	if l.end != nil {
		l.end.Next = &newItem
	}
	l.end = &newItem
	if l.start == nil {
		l.start = l.end
	}
	l.len++
	return l.end
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.start = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.end = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.start {
		return
	}

	// Remove item from its current position
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.end = i.Prev
	}

	// Move to front
	i.Next = l.start
	i.Prev = nil
	l.start.Prev = i
	l.start = i
}
