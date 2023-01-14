package internal

import "log"

func NewA(b *B, c *C) *A {
	log.Println("creating A")
	return &A{b: b, c: c}
}

type A struct {
	b *B
	c *C
}

func NewB() *B {
	log.Println("creating B")
	return &B{}
}

type B struct {
}

func NewC(d *D) *C {
	log.Println("creating C")
	return &C{d: d}
}

type C struct {
	d *D
}

func NewD() *D {
	log.Println("creating D")
	return &D{}
}

type D struct {
}
