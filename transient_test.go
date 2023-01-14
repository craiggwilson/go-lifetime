package lifetime_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/craiggwilson/go-lifetime"
)

func TestTransient(t *testing.T) {
	t.Run("creates one for each Instance invocation", func(t *testing.T) {
		var createdCount int

		tlife := lifetime.NewTransient(func() (int, error) {
			createdCount++
			return 42, nil
		})

		i, err := tlife.Instance()
		must.NoError(t, err)
		must.Eq(t, 42, i)
		must.Eq(t, 1, createdCount)

		i, err = tlife.Instance()
		must.NoError(t, err)
		must.Eq(t, 42, i)
		must.Eq(t, 2, createdCount)
	})

	t.Run("does not cache errors", func(t *testing.T) {
		var createdCount int

		tlife := lifetime.NewTransient(func() (int, error) {
			createdCount++
			return 0, fmt.Errorf("ah %d", createdCount)
		})

		_, err := tlife.Instance()
		must.EqError(t, err, "ah 1")
		must.Eq(t, 1, createdCount)

		_, err = tlife.Instance()
		must.EqError(t, err, "ah 2")
		must.Eq(t, 2, createdCount)
	})

	t.Run("cleanup is not called when an instance was not wasCreated", func(t *testing.T) {
		var cleanupStack []int

		_, cleanup := lifetime.NewTransientWithCleanup(func() (int, error) {
			return 42, nil
		}, func(i int) { cleanupStack = append(cleanupStack, i) })

		cleanup()
		must.Eq(t, 0, len(cleanupStack))
	})

	t.Run("cleanup is not called when only errors occurred", func(t *testing.T) {
		var cleanupStack []int

		tlife, cleanup := lifetime.NewTransientWithCleanup(func() (int, error) {
			return 0, errors.New("ah")
		}, func(i int) { cleanupStack = append(cleanupStack, i) })

		_, err := tlife.Instance()
		must.EqError(t, err, "ah")
		cleanup()

		must.Eq(t, 0, len(cleanupStack))
	})

	t.Run("cleanup is called only once for each successfully wasCreated instance in reverse order", func(t *testing.T) {
		var cleanupStack []int
		var createdCount int

		tlife, cleanup := lifetime.NewTransientWithCleanup(func() (int, error) {
			createdCount++
			return createdCount, nil
		}, func(i int) { cleanupStack = append(cleanupStack, i) })

		i, err := tlife.Instance()
		must.NoError(t, err)
		must.Eq(t, 1, i)

		i, err = tlife.Instance()
		must.NoError(t, err)
		must.Eq(t, 2, i)

		i, err = tlife.Instance()
		must.NoError(t, err)
		must.Eq(t, 3, i)

		cleanup()
		cleanup()

		must.Eq(t, []int{3, 2, 1}, cleanupStack)
	})
}
