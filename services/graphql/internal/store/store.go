package store

import (
	"sync"
	"time"
)

type Task struct {
	ID          string
	Title       string
	Description *string
	Done        bool
	CreatedAt   time.Time
}

type Store struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewStore() *Store {
	return &Store{
		tasks: make(map[string]Task),
	}
}

func (s *Store) Create(id, title string, description *string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now(),
	}
	s.tasks[id] = task
	return task
}

func (s *Store) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *Store) GetByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) Update(id string, title *string, description *string, done *bool) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	if title != nil {
		t.Title = *title
	}
	if description != nil {
		t.Description = description
	}
	if done != nil {
		t.Done = *done
	}
	s.tasks[id] = t
	return t, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if ok {
		delete(s.tasks, id)
	}
	return ok
}
