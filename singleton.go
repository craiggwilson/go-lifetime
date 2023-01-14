package lifetime

import (
	"sync"

	"github.com/craiggwilson/go-lifetime/internal"
)

// NewSingleton creates a Singleton lifetime.
func NewSingleton[T any](create func() (T, error)) *Singleton[T] {
	s, _ := NewSingletonWithCleanup(create, nil)
	return s
}

// NewSingletonWithCleanup creates a Singleton lifetime and returns a cleanup function that should be called
// when the instance is no longer needed.
func NewSingletonWithCleanup[T any](create func() (T, error), cleanup func(T)) (*Singleton[T], func()) {
	s := &Singleton[T]{
		create:  create,
		cleanup: cleanup,
	}

	return s, s.Cleanup
}

// Singleton creates and holds a single instance of T, returning the same instance for each invocation of Instance.
type Singleton[T any] struct {
	createdOnce sync.Once
	cleanupOnce sync.Once

	create  func() (T, error)
	cleanup func(T)

	wasCreated internal.AtomicBool
	instance   T
	err        error
}

// MustInstance calls Instance and panics if an error occurs.
func (s *Singleton[T]) MustInstance() T {
	inst, err := s.Instance()
	if err != nil {
		panic(err)
	}

	return inst
}

// Instance creates an instance of T and returns an err if one could not be created.
func (s *Singleton[T]) Instance() (T, error) {
	s.createdOnce.Do(func() {
		s.instance, s.err = s.create()
		s.wasCreated.Swap(s.err == nil)
	})

	return s.instance, s.err
}

// Cleanup cleans up the created instance if a cleanup function was provided.
func (s *Singleton[T]) Cleanup() {
	if s.cleanup != nil && s.wasCreated.Load() {
		s.cleanupOnce.Do(func() {
			s.cleanup(s.instance)
		})
	}
}
