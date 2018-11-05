package main

import "fmt"

// This tenet will match around 28 things.
// 'NamespaceAll.*List'

// | code context                   | dir              |
// |--------------------------------+------------------|
// | c.serviceClient.Services()     | ./pkg/registry   |
// | jm.kubeClient.BatchV1().Jobs() | ./pkg/controller |
// | kubeClient.CoreV1().Pods()     | ./pkg/controller |
// | fakeKubeClient.Core().Pods()   | ./pkg/controller |
// | o.client.CoreV1().Pods()       | ./pkg/kubectl    |
// | tc.Core().ServiceAccounts()    | ./pkg/client     |
// | tc.Core().Pods()               | ./pkg/client     |

// Check this function is an object of a controller
// No example of Reconcile anywhere to be seen -- would have to check git pickaxe
func (c *ControllerX) Reconcile() {
	// X here is arbitrary
	// I should not match what is in X's parameters, nor List's
	// metav1 is being used at the moment
	items, err := c.kubeClient.X(v1.NamespaceAll).List(&v1.ListOptions{})
}
