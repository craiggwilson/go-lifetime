package lifetime_test

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/craiggwilson/go-lifetime"
)

func TestSingleton(t *testing.T) {
	t.Run("only creates one", func(t *testing.T) {
		var createdCount int

		slife := lifetime.NewSingleton(func() (int, error) {
			createdCount++
			return 42, nil
		})

		i, err := slife.Instance()
		must.NoError(t, err)
		must.Eq(t, 42, i)
		must.Eq(t, 1, createdCount)

		i, err = slife.Instance()
		must.NoError(t, err)
		must.Eq(t, 42, i)
		must.Eq(t, 1, createdCount)
	})

	t.Run("returns errors", func(t *testing.T) {
		var createdCount int

		slife := lifetime.NewSingleton(func() (int, error) {
			createdCount++
			return 0, errors.New("ah")
		})

		_, err := slife.Instance()
		must.EqError(t, err, "ah")
		must.Eq(t, 1, createdCount)

		_, err = slife.Instance()
		must.EqError(t, err, "ah")
		must.Eq(t, 1, createdCount)
	})

	t.Run("cleanup is not called when an instance was not wasCreated", func(t *testing.T) {
		var cleanupCount int

		_, cleanup := lifetime.NewSingletonWithCleanup(func() (int, error) {
			return 42, nil
		}, func(int) { cleanupCount++ })

		cleanup()
		must.Eq(t, 0, cleanupCount)
	})

	t.Run("cleanup is not called when an error occurred", func(t *testing.T) {
		var cleanupCount int

		slife, cleanup := lifetime.NewSingletonWithCleanup(func() (int, error) {
			return 0, errors.New("ah")
		}, func(int) { cleanupCount++ })

		_, err := slife.Instance()
		must.EqError(t, err, "ah")
		cleanup()

		must.Eq(t, 0, cleanupCount)
	})

	t.Run("cleanup is only called once", func(t *testing.T) {
		var cleanupCount int

		slife, cleanup := lifetime.NewSingletonWithCleanup(func() (int, error) {
			return 42, nil
		}, func(int) { cleanupCount++ })

		i, err := slife.Instance()
		must.NoError(t, err)
		must.Eq(t, 42, i)

		cleanup()
		cleanup()

		must.Eq(t, 1, cleanupCount)
	})
}
