package datastructures

import "sort"

type Array[Stored comparable] interface {
	Get() []Stored
	Add(elements ...Stored)
	Set(index int, e Stored)
	Remove(indices ...int)
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

func (s *array[Stored]) Set(index int, e Stored) {
	if diffSize := index - len(s.Data) + 1; diffSize > 0 {
		diff := make([]Stored, diffSize)
		s.Add(diff...)
	}
	if s.Data[index] != e {
		s.Data[index] = e
	}
}

func (s *array[Stored]) Remove(indices ...int) {
	sort.Slice(indices, func(i, j int) bool { return indices[i] > indices[j] })
	for _, index := range indices {
		if index >= len(s.Data) {
			continue
		}
		s.Data[index] = s.Data[len(s.Data)-1]
		s.Data = s.Data[:len(s.Data)-1]
	}
}
