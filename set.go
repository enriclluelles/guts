package main

type Set map[interface{}]bool

func NewSet() *Set {
	s := make(Set)
	return &s
}

func (s *Set) Add(elem interface{}) {
	_, ok := (*s)[elem]

	if !ok {
		(*s)[elem] = true
	}
}
