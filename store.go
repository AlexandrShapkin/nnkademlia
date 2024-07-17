package nnkademlia

import (
	"fmt"
	"math/big"
)

type StoreItem struct {
	value string
	origin *big.Int
}

func NewStoreItem(value string, origin *big.Int) *StoreItem {
	return &StoreItem{
		value: value,
		origin: origin,
	}
}

func (si *StoreItem) String() string {
	return fmt.Sprintf("v: %s, o: %s", si.value, si.origin.Text(16))
}

type Store struct {
	store map[string]*StoreItem
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]*StoreItem),
	}
}

func (s *Store) String() string {
	if len(s.store) == 0 {
		return ""
	}
	store := ""
	for k, v := range s.store {
		store += fmt.Sprintf("%s: [%s], ", k, v.String())
	}
	return store[:len(store) - 2]
}

func (s *Store) Add(key string, value string, origin *big.Int) bool {
	flag := true
	if _, isRewrite := s.store[key]; isRewrite {
		flag = false
	}
	s.store[key] = NewStoreItem(
		value,
		origin,
	)

	return flag
}

func (s *Store) Find(key string) (string, bool) {
	item, ok := s.store[key]
	if !ok {
		return "", false
	}
	return item.value, true
}