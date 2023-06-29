package aggservice

import (
	"fmt"

	"github.com/mhg14/toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

type Storer interface {
	Insert(types.Distance) error
	Get(id int) (float64, error)
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	distance, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("couldn't find the aggregated distance for obu id %d", id)
	}
	return distance, nil
}
