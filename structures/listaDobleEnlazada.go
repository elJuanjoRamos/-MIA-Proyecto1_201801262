package structures

type NODE struct {
	// Next y prev son los punteros en la lista doble.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next PARTITION of the last
	// list PARTITION (l.Back()) and the previous PARTITION of the first list
	// PARTITION (l.Front()).
	next, prev *NODE

	// The list to which this PARTITION belongs.
	list *List

	// El valor, va a guardar una particion.
	Value PARTITION
}

// List representa una lista doblemente enlazada.
// El valor cero de List es una lista vacía.
type List struct {
	root NODE
	len  int
}

// Inicializa o limpia la lista l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// Retorna una nueva lista inicilizada
func New() *List { return new(List).Init() }

// Retorna la cantidad de PARTITIONos en la lista l.
func (l *List) Len() int { return l.len }

// Front devuelve el primer elemento de la lista o nil si la lista está vacía.
func (l *List) Front() *NODE {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back devuelve el ultimo elmeento de la lista o nil si la lista esta vacia.
func (l *List) Back() *NODE {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert  inserta e después de at, incrementa l.len y devuelve e.
func (l *List) insert(e, at *NODE) *NODE {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&PARTITION{Value: v}, at).
func (l *List) insertValue(v PARTITION, at *NODE) *NODE {
	return l.insert(&NODE{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *List) remove(e *NODE) *NODE {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// move moves e to next to at and returns e.
func (l *List) move(e, at *NODE) *NODE {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e

	return e
}

// Remove removes e from l if e is an PARTITION of list l.
// It returns the PARTITION value e.Value.
// The PARTITION must not be nil.
func (l *List) Remove(e *NODE) PARTITION {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero PARTITION) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new PARTITION e with value v at the front of list l and returns e.
func (l *List) PushFront(v PARTITION) *NODE {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new PARTITION e with value v at the back of list l and returns e.
func (l *List) PushBack(v PARTITION) *NODE {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new PARTITION e with value v immediately before mark and returns e.
// If mark is not an PARTITION of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertBefore(v PARTITION, mark *NODE) *NODE {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new PARTITION e with value v immediately after mark and returns e.
// If mark is not an PARTITION of l, the list is not modified.
// The mark must not be nil.
func (l *List) InsertAfter(v PARTITION, mark *NODE) *NODE {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves PARTITION e to the front of list l.
// If e is not an PARTITION of l, the list is not modified.
// The PARTITION must not be nil.
func (l *List) MoveToFront(e *NODE) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves PARTITION e to the back of list l.
// If e is not an PARTITION of l, the list is not modified.
// The PARTITION must not be nil.
func (l *List) MoveToBack(e *NODE) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves PARTITION e to its new position before mark.
// If e or mark is not an PARTITION of l, or e == mark, the list is not modified.
// The PARTITION and mark must not be nil.
func (l *List) MoveBefore(e, mark *NODE) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves PARTITION e to its new position after mark.
// If e or mark is not an PARTITION of l, or e == mark, the list is not modified.
// The PARTITION and mark must not be nil.
func (l *List) MoveAfter(e, mark *NODE) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}
