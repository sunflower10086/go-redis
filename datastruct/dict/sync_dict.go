package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{
		m: sync.Map{},
	}
}

func (s *SyncDict) Get(key string) (val any, exists bool) {
	return s.m.Load(key)
}

func (s *SyncDict) Len() int {
	size := 0
	s.m.Range(func(key, value any) bool {
		size++
		return true
	})

	return size
}

func (s *SyncDict) Put(key string, val any) (result int) {
	_, existed := s.m.Load(key)
	s.m.Store(key, val)
	if existed {
		return 0
	}

	return 1
}

func (s *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, existed := s.m.Load(key)
	if existed {
		return
	}
	s.m.Store(key, val)

	return 1
}

func (s *SyncDict) PutIfExists(key string, val any) (result int) {
	_, existed := s.m.Load(key)
	if existed {
		s.m.Store(key, val)
		return 1
	}

	return 0
}

func (s *SyncDict) Remove(key string) (result int) {
	_, existed := s.m.Load(key)
	if existed {
		s.m.Delete(key)
		return 1
	}

	return 0
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncDict) Keys() []string {
	keys := make([]string, 0, s.Len())
	s.m.Range(func(key, value any) bool {
		keys = append(keys, key.(string))
		return true
	})

	return keys
}

func (s *SyncDict) RandomKeys(limit int) []string {
	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		s.m.Range(func(key, value any) bool {
			keys = append(keys, key.(string))
			return false
		})
	}

	return keys
}

func (s *SyncDict) RandomDistinckKeys(limit int) []string {
	keys := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		idx := i
		s.m.Range(func(key, value any) bool {
			keys = append(keys, key.(string))
			if idx == limit {
				return false
			}
			return true
		})
	}

	return keys
}

func (s *SyncDict) Clear() {
	*s = *NewSyncDict()
}
