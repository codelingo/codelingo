package main

import (
	"fmt"
)

func aFunc(ctx context.Context, a int) {
	// Do something
}

func bFunc(b int, ctx context.Context, a int) {
	// Do something
}

func cFunc(c int, b int, ctx context.Context, a int) {
	// Do something
}

// Don't catch them in structs
type IsTypeOfParams struct {
	// Value that needs to be resolve.
	// Use this to decide which GraphQLObject this value maps to.
	Value interface{}

	// Info is a collection of information about the current execution state.
	Info ResolveInfo

	// Context argument is a context value that is provided to every resolve function within an execution.
	// It is commonly
	// used to represent an authenticated user, or request-specific caches.
	Context context.Context
}
