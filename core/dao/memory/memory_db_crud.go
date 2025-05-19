package memory

import (
	"fmt"

	"github.com/google/uuid"
)

type Identifiable interface {
	SetId(id string)
	GetId() string
}

type MemoryDbCrud[T Identifiable] struct {
	table string
	data  map[string]*T
}

func NewMemoryDbCrud[T Identifiable](table string) *MemoryDbCrud[T] {
	return &MemoryDbCrud[T]{
		table: table,
		data:  make(map[string]*T),
	}
}

func (dao *MemoryDbCrud[T]) Create(data T) (T, error) {
	id := uuid.NewString()
	dao.data[id] = &data
	data.SetId(id)
	return data, nil
}

func (dao *MemoryDbCrud[T]) Update(id string, data T) (T, error) {
	if dao.data[id] == nil {
		return data, fmt.Errorf("could not update %s with id %s", dao.table, id)
	}
	dao.data[id] = &data
	return data, nil
}

func (dao *MemoryDbCrud[T]) ReadById(id string) T {
	return *dao.data[id]
}

func (dao *MemoryDbCrud[T]) ReadAll() []T {
	var result []T
	for _, item := range dao.data {
		result = append(result, *item)
	}
	return result
}
