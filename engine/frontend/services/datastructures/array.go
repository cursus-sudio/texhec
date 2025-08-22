package datastructures

import "sort"

type Array[Stored comparable] interface {
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
}

type array[Stored comparable] struct {
	data []Stored
}

func NewArray[Stored comparable]() Array[Stored] {
	return &array[Stored]{}
}

func (s *array[Stored]) Get() []Stored { return s.data }

func (s *array[Stored]) Add(elements ...Stored) {
	s.data = append(s.data, elements...)
}

func (s *array[Stored]) Set(index int, e Stored) error {
	if len(s.data) <= index {
		return ErrOutOfBounds
	}
	if s.data[index] == e {
		return nil
	}
	s.data[index] = e

	return nil
}

func (s *array[Stored]) Remove(indices ...int) error {
	for _, index := range indices {
		if index >= len(s.data) {
			return ErrOutOfBounds
		}
	}

	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		s.data[index] = s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
	}
	return nil
}
