package goid

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/leromarinvit/typist"
)

func TestGoid(t *testing.T) {
	done := make(chan int)
	compareGoid(100, t, done)
	<-done
}

func compareGoid(count int, t *testing.T, done chan<- int) {
	goid := Goid()
	stack := goidFromStack()
	if goid != stack {
		t.Fatalf("Goid() returned %d, stack trace said %d\n", goid, stack)
	}
	if count == 0 {
		done <- 1
	} else {
		go compareGoid(count-1, t, done)
	}
}

func goidFromStack() (id int64) {
	buf := make([]byte, 20)
	runtime.Stack(buf, false)
	num := string(buf[len("goroutine "):])
	num = num[:strings.IndexByte(num, ' ')]
	id, _ = strconv.ParseInt(num, 10, 64)
	return id
}

func BenchmarkGoid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Goid()
	}
}

func TestRuntimeG(t *testing.T) {
	pg, err := typist.TypeByString("*runtime.g")
	if err != nil {
		t.Fatal("*runtime.g not found")
	}
	if pg.Kind() != reflect.Ptr {
		t.Fatal("*runtime.g isn't reflect.Ptr but", pg.Kind())
	}
	g := pg.Elem()
	if g.Kind() != reflect.Struct {
		t.Fatal("runtime.g isn't reflect.Struct but", g.Kind())
	}
	goid, ok := g.FieldByName("goid")
	if !ok {
		t.Fatal("runtime.g has no goid field")
	}
	if goid.Type.Kind() != reflect.Int64 {
		t.Fatal("runtime.g.goid isn't int64 but", goid.Type.String())
	}
	_, ok = g.FieldByName("sched")
	if !ok {
		t.Fatal("runtime.g.sched has no sched field")
	}
}
