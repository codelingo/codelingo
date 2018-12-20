package main

import (
	"fmt"
)

func F(ctx context.Context, a int)               {}
func G(b int, ctx context.Context, a int)        {}
func H(c int, b int, ctx context.Context, a int) {}


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