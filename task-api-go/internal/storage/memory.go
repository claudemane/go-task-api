package storage

import (
    "errors"
    "sync"

    "taskapi/internal/models"
)

var (
    ErrNotFound = errors.New("task not found")
)

// MemoryStore stores tasks in-memory and is safe for concurrent access.
type MemoryStore struct {
    mu     sync.Mutex
    nextID int
    tasks  map[int]models.Task
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        nextID: 1,
        tasks:  make(map[int]models.Task),
    }
}

func (s *MemoryStore) Create(title string) models.Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    t := models.Task{
        ID:    s.nextID,
        Title: title,
        Done:  false,
    }
    s.tasks[t.ID] = t
    s.nextID++
    return t
}

func (s *MemoryStore) Get(id int) (models.Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    t, ok := s.tasks[id]
    if !ok {
        return models.Task{}, ErrNotFound
    }
    return t, nil
}

func (s *MemoryStore) List() []models.Task {
    s.mu.Lock()
    defer s.mu.Unlock()

    out := make([]models.Task, 0, len(s.tasks))
    for _, t := range s.tasks {
        out = append(out, t)
    }
    return out
}

func (s *MemoryStore) UpdateDone(id int, done bool) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    t, ok := s.tasks[id]
    if !ok {
        return ErrNotFound
    }
    t.Done = done
    s.tasks[id] = t
    return nil
}
