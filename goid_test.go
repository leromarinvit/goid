package goid

import (
	"runtime"
	"strconv"
	"strings"
	"testing"
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
		t.Errorf("Goid() returned %d, stack trace said %d\n", goid, stack)
		t.Fail()
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
