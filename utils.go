package graphql

type void struct{}

type set struct {
	elems map[string]void
}

func newSet(elems ...string) *set {
	st := &set{
		elems: make(map[string]void),
	}
	for _, elem := range elems {
		st.add(elem)
	}
	return st
}

func (s *set) add(value string) {
	s.elems[value] = void{}
}

func (s *set) has(value string) bool {
	_, ok := s.elems[value]
	return ok
}

func (s *set) elements() []string {
	elems := []string{}
	for key := range s.elems {
		elems = append(elems, key)
	}
	return elems
}
