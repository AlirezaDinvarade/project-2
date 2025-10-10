package main

import types "tolling/Types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Insert(d types.Distance) error {
	return nil
}
