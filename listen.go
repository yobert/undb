package undb

func (store *Store) Listen(c chan Op) {
	if store.listeners == nil {
		store.listeners = make(map[chan Op]struct{})
	}
	store.listeners[c] = struct{}{}
}

func (store *Store) Unlisten(c chan Op) {
	if store.listeners == nil {
		return
	}
	delete(store.listeners, c)
}

func (store *Store) Emit(op *Op, source string) {
	o := op.Copy()
	s := store

	o.changesource = source

	path := s.Id

	for s != nil {

		if s.parent == nil {
			s.dirty = true
		}

		o.Path = path
		for c := range s.listeners {
			c <- o
		}
		s = s.parent
		if s != nil {
			path = s.Id + "." + path
		}
	}
}
