package main

import (
	"context"
	"fmt"
	"time"
)

func foo(ctx context.Context, c chan string) {
	go func() {
		for i := 0; i < 5; i++ {
			select {
			case c <- "message":
			case <-ctx.Done(): // Accepted
				return
			}
		}
		close(c)
	}()
}

func bar(ctx context.Context, c1 chan string) {
	go func() {
		for i := 0; i < 5; i++ {
			c1 <- "message1" // Issue
		}
		close(c1)
	}()
}

func main() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	c := make(chan string)

	foo(ctx, c)

	for msg := range c {
		fmt.Println(msg)
	}

	c1 := make(chan string)

	bar(ctx, c1)

	for msg := range c1 {
		fmt.Println(msg)
	}
}
