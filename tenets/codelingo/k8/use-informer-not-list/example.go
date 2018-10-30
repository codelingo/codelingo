package main

import "fmt"

func (c *ControllerX) Reconcile() {
  items, err := c.kubeClient.X(v1.NamespaceAll).List(&v1.ListOptions{})
  if err != nil {
    ...
  }
  for _, item := range items {
    ...
  }
}
 
func (c *ControllerX) Run() {
  wait.Until(c.Reconcile, c.Period, wait.NeverStop)
  ...
}

func main() {
}

// Heuristic:
// Look for List(). Probably not possible to replace easily with Informer. grep for examples.
