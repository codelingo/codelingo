package sync

import "fmt"

type syncOp interface {
	area() float64
}
