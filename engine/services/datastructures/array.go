package datastructures

import "sort"

type Array[Stored comparable] interface {
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored) error
	Remove(indices ...int) error
}

type array[Stored comparable] struct {
	Data []Stored
}

func NewArray[Stored comparable]() Array[Stored] {
	return &array[Stored]{}
}

func (s *array[Stored]) Get() []Stored { return s.Data }

func (s *array[Stored]) Add(elements ...Stored) {
	s.Data = append(s.Data, elements...)
}

func (s *array[Stored]) Set(index int, e Stored) error {
	if len(s.Data) <= index {
		return ErrOutOfBounds
	}
	if s.Data[index] == e {
		return nil
	}
	s.Data[index] = e

	return nil
}

func (s *array[Stored]) Remove(indices ...int) error {
	for _, index := range indices {
		if index >= len(s.Data) {
			return ErrOutOfBounds
		}
	}

	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		s.Data[index] = s.Data[len(s.Data)-1]
		s.Data = s.Data[:len(s.Data)-1]
	}
	return nil
}
