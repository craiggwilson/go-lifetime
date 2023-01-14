package lifetime

// Lifetime manages the lifetime of instances of T.
type Lifetime[T any] interface {
	// MustInstance calls Instance and panics if an error occurs.
	MustInstance() T
	// Instance creates an instance of T and returns an err if one could not be created.
	Instance() (T, error)
}
