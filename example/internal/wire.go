//go:build wireinject
// +build wireinject

package internal

import (
	"fmt"
	"log"

	"github.com/craiggwilson/go-lifetime"

	"github.com/google/wire"
)

func makeA(cfg *Config, blife lifetime.Lifetime[*B], clife lifetime.Lifetime[*C]) (lifetime.Lifetime[*A], func()) {
	return lifetime.NewSingletonWithCleanup(func() (*A, error) {
		var b *B
		var err error
		if cfg.UseB {
			b, err := blife.Instance()
			if err != nil {
				return nil, fmt.Errorf("creating B: %w", b)
			}
		}

		var c *C
		if cfg.UseC {
			c, err = clife.Instance()
			if err != nil {
				return nil, fmt.Errorf("creating C: %w", c)
			}
		}

		return NewA(b, c), nil
	}, func(a *A) { log.Println("cleaning up A") })
}

func makeB(cfg *Config) lifetime.Lifetime[*B] {
	return lifetime.NewSingleton(func() (*B, error) { return NewB(), nil })
}

func makeC(cfg *Config, dlife lifetime.Lifetime[*D]) (lifetime.Lifetime[*C], func()) {
	return lifetime.NewSingletonWithCleanup(func() (*C, error) {
		d, err := dlife.Instance()
		if err != nil {
			return nil, fmt.Errorf("creating D: %w", d)
		}

		return NewC(d), nil
	}, func(c *C) { log.Println("cleaning up C") })
}

func makeD(cfg *Config) (lifetime.Lifetime[*D], func()) {
	return lifetime.NewSingletonWithCleanup(func() (*D, error) {
		return NewD(), nil
	}, func(d *D) { log.Println("cleaning up D") })
}

func New(cfg *Config) (lifetime.Lifetime[*A], func()) {
	panic(wire.Build(
		makeA,
		makeB,
		makeC,
		makeD,
	))
}
