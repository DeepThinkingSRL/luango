package entity

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type ID int

type Entity struct {
	ID       ID
	Name     string
	Position struct {
		X, Y float64
	}
	Sprite *ebiten.Image
}

type Manager struct {
	entities map[ID]*Entity
	nextID   ID
	lock     sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		entities: make(map[ID]*Entity),
		nextID:   1,
	}
}

func (em *Manager) CreateEntity(name string, sprite *ebiten.Image) *Entity {
	em.lock.Lock()
	defer em.lock.Unlock()

	e := &Entity{
		ID:     em.nextID,
		Name:   name,
		Sprite: sprite,
	}
	em.entities[e.ID] = e
	em.nextID++
	return e
}

func (em *Manager) GetEntity(id ID) (*Entity, bool) {
	em.lock.Lock()
	defer em.lock.Unlock()
	e, ok := em.entities[id]
	return e, ok
}

func (em *Manager) GetAllEntities() map[ID]*Entity {
	em.lock.Lock()
	defer em.lock.Unlock()
	// Return a copy to avoid race conditions
	result := make(map[ID]*Entity)
	for k, v := range em.entities {
		result[k] = v
	}
	return result
}

func (em *Manager) GetEntitiesSlice() []*Entity {
	em.lock.Lock()
	defer em.lock.Unlock()
	
	entities := make([]*Entity, 0, len(em.entities))
	for _, e := range em.entities {
		entities = append(entities, e)
	}
	return entities
}

func (em *Manager) RemoveEntity(id ID) bool {
	em.lock.Lock()
	defer em.lock.Unlock()
	
	if _, exists := em.entities[id]; exists {
		delete(em.entities, id)
		return true
	}
	return false
}

func (em *Manager) Count() int {
	em.lock.Lock()
	defer em.lock.Unlock()
	return len(em.entities)
}
