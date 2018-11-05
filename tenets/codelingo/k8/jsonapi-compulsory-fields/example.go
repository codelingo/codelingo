package main

// Heuristic:
// ‾‾‾‾‾‾‾‾‾‾
// this used to live inside of staging. but i have removed that when testing kubernetes
// go files deep inside the 'k8s.io/api/' folder. find .go files that import metav1 k8s.io/apimachinery/pkg/apis/meta/v1 and define a bunch of structs. detecting import and struct may not be necessary. We want to find files that do not have kind and apiVersion

import "fmt"

func main() {
}
