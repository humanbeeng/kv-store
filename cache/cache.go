package cache

import (
	"fmt"
	"sync"
)

var StringStore *KVStore[string, string]

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
}

func (s *KVStore[K, V]) Has(key K) bool {
	_, found := s.store[key]
	return found
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		store: make(map[K]V),
	}
}

// Time to make our KVStore a Storer
func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
	return nil
}

func (s *KVStore[K, V]) Update(key K, val V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.store[key]
	if !found {
		return fmt.Errorf("key %v not found", key)
	}
	s.store[key] = val
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	val, found := s.store[key]

	if !found {
		return val, fmt.Errorf("key %v does not exists", key)
	}
	return val, nil
}

func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.Lock()

	val, found := s.store[key]
	if !found {
		return val, fmt.Errorf("key %v not found", key)
	}
	delete(s.store, key)
	return val, nil
}

func init() {
	StringStore = NewKVStore[string, string]()
}
