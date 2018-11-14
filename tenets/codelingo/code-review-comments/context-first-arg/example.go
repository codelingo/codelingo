// Package main used for tenet testing
package main

import (
	"fmt"
)

func a(ctx context.Context, a int)               {}

func b(b int, ctx context.Context, a int)        {}

func c(c int, b int, ctx context.Context, a int) {}
