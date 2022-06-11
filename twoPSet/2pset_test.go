package twoPSet

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	set := New()
	assert.NotNil(t, set)
	assert.Empty(t, set.GetAddSet())
	assert.Empty(t, set.GetRemoveSet())
}

func TestT_Add(t *testing.T) {
	set := New()
	id := uuid.New()
	payload := []byte("hello")
	set.Add(id, payload)

	addSet := set.GetAddSet()
	assert.NotEmpty(t, addSet)
	assert.Contains(t, addSet, id)
	assert.Equal(t, payload, addSet[id].Payload)
}

func TestT_Remove(t *testing.T) {
	set := New()
	id := uuid.New()
	payload := []byte("hello")
	set.Add(id, payload)

	err := set.Remove(id)
	assert.NoError(t, err)

	removeSet := set.GetRemoveSet()
	assert.Contains(t, removeSet, id)
	assert.Equal(t, payload, removeSet[id].Payload)
}

func TestT_Remove_ElementDoesntExist(t *testing.T) {
	set := New()
	id := uuid.New()

	err := set.Remove(id)
	assert.EqualError(t, err, "element does not exist")

	removeSet := set.GetRemoveSet()
	assert.NotContains(t, removeSet, id)
}

func TestT_Merge(t *testing.T) {
	set1 := New()
	id1 := uuid.New()
	payload1 := []byte("hello")
	set1.Add(id1, payload1)
	err := set1.Remove(id1)
	assert.NoError(t, err)

	set2 := New()
	id2 := uuid.New()
	payload2 := []byte("world")
	set2.Add(id2, payload2)
	err = set2.Remove(id2)
	assert.NoError(t, err)

	set1.Merge(set2)

	assert.Contains(t, set1.AddSet, id1)
	assert.Contains(t, set1.AddSet, id2)
	assert.Contains(t, set1.RemoveSet, id1)
	assert.Contains(t, set1.RemoveSet, id2)
}

func TestMerge(t *testing.T) {
	setA := New()
	setB := New()

	id1 := uuid.New()
	payload1 := []byte("hello")
	setA.Add(id1, payload1)

	id2 := uuid.New()
	payload2 := []byte("world")
	setB.Add(id2, payload2)

	setMerged := Merge(setA.AddSet, setB.AddSet)
	assert.Contains(t, setMerged, id1)
	assert.Contains(t, setMerged, id2)
}

func TestMerge_Conflict(t *testing.T) {
	setA := New()
	setB := New()

	id := uuid.New()
	payload1 := []byte("hello")
	setA.Add(id, payload1)

	time.Sleep(1)
	payload2 := []byte("world")
	setB.Add(id, payload2)

	setMerged := Merge(setA.AddSet, setB.AddSet)
	assert.Contains(t, setMerged, id)
	assert.Equal(t, payload2, setMerged[id].Payload)
}
