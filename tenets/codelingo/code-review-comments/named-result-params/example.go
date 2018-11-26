//Package main is an example package
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}

// Parent1 is an example function
func (n *Node) Parent1() (node *Node) {}

// Parent2 is an example function
func (n *Node) Parent2() (node Node, err error) {}

func (h *ClientHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {}
