package space

import (
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Space interface {
	Key() string
	Get(key string) (any, bool)
	Set(key string, value any)
	Has(key string) bool
	Remove(key string)
	Clear()
}

// region spaceData
type spaceData struct {
	key    string
	space  cmap.ConcurrentMap[string, any]
	parent Space
}

func (s *spaceData) Key() string {
	return s.key
}

func (s *spaceData) Get(key string) (any, bool) {
	val, has := s.space.Get(key)
	if has {
		return val, has
	} else if s.parent != nil {
		return s.parent.Get(key)
	}
	return nil, false
}

func (s *spaceData) Set(key string, val any) {
	s.space.Set(key, val)
}

func (s *spaceData) selfHas(key string) bool {
	return s.space.Has(key)
}

func (s *spaceData) Has(key string) bool {
	has := s.selfHas(key)
	if !has && s.parent != nil {
		has = s.parent.Has(key)
	}
	return has
}

func (s *spaceData) Remove(key string) {
	if s.selfHas(key) {
		s.space.Remove(key)
	}
}

func (s *spaceData) Clear() {
	s.space.Clear()
}

// endregion

func newSpace(key string) Space {
	return newSpaceWithParent(key, nil)
}

func newSpaceWithParent(key string, parent Space) Space {
	space := &spaceData{
		key:    key,
		space:  cmap.New[any](),
		parent: parent,
	}
	return space
}
