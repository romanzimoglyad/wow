package storage

import (
	"errors"
	"sync"
)

var NotFound = errors.New("NotFound")

type Storage struct {
	*sync.RWMutex
	m map[string]struct{}
}

func New() *Storage {
	return &Storage{
		RWMutex: &sync.RWMutex{},
		m:       make(map[string]struct{}),
	}
}

func (h *Storage) Add(hash string) error {
	h.Lock()
	defer h.Unlock()
	h.m[hash] = struct{}{}
	return nil
}

func (h *Storage) Get(hash string) error {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.m[hash]; ok {
		return nil
	}

	return NotFound
}
func (h *Storage) Delete(hash string) error {
	if err := h.Get(hash); err != nil {
		return err
	}
	h.Lock()
	defer h.Unlock()
	delete(h.m, hash)
	return nil
}
