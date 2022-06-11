package twoPSet

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type TwoPSet interface {
	GetAddSet() Set
	GetRemoveSet() Set
	Add(uuid.UUID, interface{})
	Remove(uuid.UUID) error
	Merge(TwoPSet)
}

type OP struct {
	Payload   interface{}
	Timestamp time.Time
}

type Set map[uuid.UUID]OP

type T struct {
	AddSet    Set
	RemoveSet Set
}

func New() *T {
	return &T{
		AddSet:    make(Set, 0),
		RemoveSet: make(Set, 0),
	}
}

var _ TwoPSet = &T{}

func (t *T) GetAddSet() Set {
	return t.AddSet
}

func (t *T) GetRemoveSet() Set {
	return t.RemoveSet
}

func (t *T) Add(id uuid.UUID, payload interface{}) {
	t.AddSet[id] = OP{
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

func (t *T) Remove(id uuid.UUID) error {
	if _, ok := t.AddSet[id]; !ok {
		return errors.New("element does not exist")
	}

	t.RemoveSet[id] = OP{
		Timestamp: time.Now(),
		Payload:   t.AddSet[id].Payload,
	}
	return nil
}

func (t *T) Merge(set TwoPSet) {
	t.AddSet = Merge(t.AddSet, set.GetAddSet())
	t.RemoveSet = Merge(t.RemoveSet, set.GetRemoveSet())
}

func Merge(setA, setB Set) Set {
	for k, v := range setB {
		n, ok := setA[k]
		if ok {
			if n.Timestamp.Before(v.Timestamp) {
				setA[k] = v
			}
		} else {
			setA[k] = v
		}
	}

	return setA
}
