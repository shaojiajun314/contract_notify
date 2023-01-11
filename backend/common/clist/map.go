package clist

import (
	"container/heap"
	"sort"
)

// SortedMap is a uint64->data hash map with a heap based index to allow
// iterating over the contents in a index-incrementing way.
type SortedMap struct {
	items map[uint64]interface{} // Hash map storing the data
	index *Heap                  // Heap of nonces of all the stored datas (non-strict mode)
	cache []interface{}          // Cache of the datas already sorted
}

// NewSortedMap creates a new uint64-sorted map.
func NewSortedMap() *SortedMap {
	return &SortedMap{
		items: make(map[uint64]interface{}),
		index: new(Heap),
	}
}

// Get retrieves the current datas associated with the given index.
func (m *SortedMap) Get(index uint64) interface{} {
	return m.items[index]
}

// Put inserts a new  into the map, also updating the map's index
// index. If a  already exists with the same index, it's overwritten.
func (m *SortedMap) Put(index uint64, data interface{}) {
	if m.items[index] == nil {
		heap.Push(m.index, index)
	}
	m.items[index], m.cache = data, nil
}

// Forward removes all data from the map with a index lower than the
// provided threshold.
func (m *SortedMap) Forward(threshold uint64) []interface{} {
	var removed []interface{}

	// Pop off heap items until the threshold is reached
	for m.index.Len() > 0 && (*m.index)[0] < threshold {
		index := heap.Pop(m.index).(uint64)
		removed = append(removed, m.items[index])
		delete(m.items, index)
	}
	// If we had a cached order, shift the front
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

func (m *SortedMap) Filter(filter func(interface{}) bool) []interface{} {
	removed := m.filter(filter)
	if len(removed) > 0 {
		m.reheap()
	}
	return removed
}

func (m *SortedMap) reheap() {
	*m.index = make([]uint64, 0, len(m.items))
	for index := range m.items {
		*m.index = append(*m.index, index)
	}
	heap.Init(m.index)
	m.cache = nil
}

// filter is identical to Filter, but **does not** regenerate the heap. This method
// should only be used if followed immediately by a call to Filter or reheap()
func (m *SortedMap) filter(filter func(interface{}) bool) []interface{} {
	var removed []interface{}

	// Collect all the transactions to filter out
	for index, data := range m.items {
		if filter(data) {
			removed = append(removed, data)
			delete(m.items, index)
		}
	}
	if len(removed) > 0 {
		m.cache = nil
	}
	return removed
}

// Cap places a hard limit on the number of items, returning all data exceeding that limit.
func (m *SortedMap) Cap(threshold int) []interface{} {
	// Short circuit if the number of items is under the limit
	if len(m.items) <= threshold {
		return nil
	}
	// Otherwise gather and drop the highest index'd data
	var drops []interface{}

	sort.Sort(*m.index)
	for size := len(m.items); size > threshold; size-- {
		drops = append(drops, m.items[(*m.index)[size-1]])
		delete(m.items, (*m.index)[size-1])
	}
	*m.index = (*m.index)[:threshold]
	heap.Init(m.index)

	// If we had a cache, shift the back
	if m.cache != nil {
		m.cache = m.cache[:len(m.cache)-len(drops)]
	}
	return drops
}

// Remove deletes a  from the maintained map, returning whether the	was found.
func (m *SortedMap) Remove(index uint64) bool {
	// Short circuit if no  is present
	_, ok := m.items[index]
	if !ok {
		return false
	}
	// Otherwise delete the  and fix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == index {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, index)
	m.cache = nil

	return true
}

func (m *SortedMap) Ready(start uint64) []interface{} {
	// Short circuit if no transactions are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental transactions
	var ready []interface{}
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}
	m.cache = nil

	return ready
}

// Len returns the length of the  map.
func (m *SortedMap) Len() int {
	return len(m.items)
}

func (m *SortedMap) flatten() []interface{} {
	// If the sorting was not cached yet, create and cache it
	if m.cache == nil {
		sortDatas := make([]sortData, 0, len(m.items))
		for index, data := range m.items {
			sortDatas = append(sortDatas, sortData{
				index: index,
				data:  data,
			})
		}
		sort.Sort(byIndex(sortDatas))
		for _, data := range sortDatas {
			m.cache = append(m.cache, data)
		}
	}
	return m.cache
}

func (m *SortedMap) Flatten() []interface{} {
	cache := m.flatten()
	datas := make([]interface{}, len(cache))
	copy(datas, cache)
	return datas
}

func (m *SortedMap) LastElement() interface{} {
	cache := m.flatten()
	return cache[len(cache)-1]
}

type sortData struct {
	index uint64
	data  interface{}
}

type byIndex []sortData

func (s byIndex) Len() int           { return len(s) }
func (s byIndex) Less(i, j int) bool { return s[i].index < s[j].index }
func (s byIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
