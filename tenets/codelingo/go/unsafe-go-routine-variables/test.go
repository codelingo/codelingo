package main

func safe() <-chan int {
	intc := make(chan int)
	go func(intc chan<- int) {
		defer close(intc)
	}(intc)
	return intc
}

func unsafe() <-chan int { // ISSUE
	intc := make(chan int)
	go func() {
		defer close(intc)
	}()
	return intc
}
