package datastructures_test

import (
	"fmt"
	"shared/services/datastructures"
	"testing"
)

func BenchmarkMapGetByIndex(b *testing.B) {
	m := map[int]string{}
	for i := 0; i < 1000; i++ {
		m[i] = fmt.Sprintf("%d", i)
	}
	b.StartTimer()
	var s string
	for i := 0; i < b.N; i++ {
		s = m[0]
	}
	if false {
		panic(s)
	}
}

func BenchmarkIndexTrackerGetByIndex(b *testing.B) {
	m := datastructures.NewSet[string]()
	for i := 0; i < 1000; i++ {
		m.Add(fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.GetStored(0)
	}
}

func BenchmarkIndexTrackerGetByValue(b *testing.B) {
	m := datastructures.NewSet[string]()
	for i := 0; i < 1000; i++ {
		m.Add(fmt.Sprintf("%d", i))
	}
	m.GetIndex("0")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.GetIndex("0")
	}
}

func BenchmarkMapIterate(b *testing.B) {
	m := map[int]string{}
	for i := 0; i < 1000; i++ {
		m[i] = fmt.Sprintf("%d", i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for range m {
		}
	}
}

func BenchmarkIndexTrackerIterate(b *testing.B) {
	m := datastructures.NewSet[string]()
	for i := 0; i < 1000; i++ {
		m.Add(fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for range m.Get() {
		}
	}
}

func BenchmarkMapAdd(b *testing.B) {
	m := map[int]string{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[i] = fmt.Sprintf("%d", i)
	}

}

func BenchmarkIndexTrackerAdd(b *testing.B) {
	m := datastructures.NewSet[string]()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Add(fmt.Sprintf("%d", i))
	}
}

func BenchmarkMapDelete(b *testing.B) {
	m := map[int]string{}
	for i := 0; i < b.N; i++ {
		m[i] = fmt.Sprintf("%d", i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		delete(m, i)
	}
}

func BenchmarkIndexTrackerDelete(b *testing.B) {
	m := datastructures.NewSet[string]()
	for i := 0; i < b.N; i++ {
		m.Add(fmt.Sprintf("%d", i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(i)
	}
}
