package main

import (
	"log"

	"github.com/craiggwilson/go-lifetime/example/internal"
)

func main() {

	cfg := &internal.Config{
		UseB: true,
	}

	alife, cleanup := internal.New(cfg)
	defer cleanup()

	_, err := alife.Instance()
	if err != nil {
		log.Panicln(err)
	}

}
