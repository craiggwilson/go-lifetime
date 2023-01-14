package lifetime

import "sync"

// NewTransient creates a Transient lifetime.
func NewTransient[T any](create func() (T, error)) *Transient[T] {
	t, _ := NewTransientWithCleanup(create, nil)
	return t
}

// NewTransientWithCleanup creates a Transient lifetime and returns a cleanup function that should be called
// when the instance is no longer needed.
func NewTransientWithCleanup[T any](create func() (T, error), cleanup func(T)) (*Transient[T], func()) {
	t := &Transient[T]{
		create:  create,
		cleanup: cleanup,
	}

	return t, t.Cleanup
}

// Transient creates a new instance of T for every call to Instance().
// When a cleanup function was provided, it also holds all the created instances in a slice and invokes the cleanup
// function in reverse instantiation order during cleanup.
type Transient[T any] struct {
	create func() (T, error)

	instances      []T
	instancesMutex sync.Mutex

	cleanup func(T)
}

// MustInstance calls Instance and panics if an error occurs.
func (t *Transient[T]) MustInstance() T {
	inst, err := t.Instance()
	if err != nil {
		panic(err)
	}

	return inst
}

// Instance creates an instance of T and returns an err if one could not be created.
func (t *Transient[T]) Instance() (T, error) {
	inst, err := t.create()
	if err == nil && t.cleanup != nil {
		t.instancesMutex.Lock()
		defer t.instancesMutex.Unlock()

		t.instances = append(t.instances, inst)
	}

	return inst, err
}

// Cleanup cleans up all the created instances if a cleanup function was provided.
func (t *Transient[T]) Cleanup() {
	if t.cleanup != nil {
		t.instancesMutex.Lock()
		defer t.instancesMutex.Unlock()

		for i := len(t.instances) - 1; i >= 0; i-- {
			t.cleanup(t.instances[i])
		}

		t.instances = nil
	}
}
