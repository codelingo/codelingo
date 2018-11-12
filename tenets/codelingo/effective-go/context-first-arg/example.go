package main

import (
	"fmt"
)

func F(ctx context.Context, a int)               {}
func F(b int, ctx context.Context, a int)        {}
func F(c int, b int, ctx context.Context, a int) {}
