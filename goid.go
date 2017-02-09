// Package goid returns the current Goroutine ID.
package goid

import (
	"reflect"
	"unsafe"

	"github.com/leromarinvit/typist"
)

//go:linkname mcall runtime.mcall
func mcall(func(uintptr))

//go:linkname gogo runtime.gogo
func gogo(uintptr)

var offsetGoid uintptr
var offsetSched uintptr

func init() {
	pg, err := typist.TypeByString("*runtime.g")
	if err != nil {
		return
	}
	if pg.Kind() != reflect.Ptr {
		return
	}
	g := pg.Elem()
	if g.Kind() != reflect.Struct {
		return
	}
	goid, ok := g.FieldByName("goid")
	if !ok {
		return
	}
	sched, ok := g.FieldByName("sched")
	if !ok {
		return
	}
	offsetGoid = goid.Offset
	offsetSched = sched.Offset
}

/*
Goid returns the current Goroutine ID.

It relies on the Go runtime's internals, so it's your own damn fault if it eats
your dog. You have been warned.

BUGS: Occasionally eats your dog
*/
func Goid() (goid int64) {
	if offsetGoid == offsetSched {
		return -1
	}
	mcall(func(g uintptr) {
		id := unsafe.Pointer(g + offsetGoid)
		goid = *(*int64)(id)
		sched := g + offsetSched
		gogo(sched)
	})
	return goid
}
